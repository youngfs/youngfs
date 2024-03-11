package server

import (
	"bytes"
	"context"
	"crypto/md5"
	"github.com/klauspost/reedsolomon"
	"github.com/youngfs/youngfs/pkg/errors"
	"github.com/youngfs/youngfs/pkg/fs/bucket"
	"github.com/youngfs/youngfs/pkg/fs/engine"
	"github.com/youngfs/youngfs/pkg/fs/entry"
	"github.com/youngfs/youngfs/pkg/fs/fullpath"
	"github.com/youngfs/youngfs/pkg/fs/meta"
	"github.com/youngfs/youngfs/pkg/util"
	"github.com/youngfs/youngfs/pkg/util/httputil"
	"github.com/youngfs/youngfs/pkg/util/mem"
	"go.uber.org/multierr"
	"io"
	"math/rand"
	"mime"
	"os"
	"path/filepath"
	"time"
)

type Server struct {
	metaStore      meta.Store
	chunkStore     engine.Engine
	endpointsCount int
}

func NewServer(meta meta.Store, chunk engine.Engine) *Server {
	cnt := 0
	endpoints, err := chunk.GetEndpoints(context.Background())
	if err == nil {
		cnt = len(endpoints)
	}

	svr := &Server{
		metaStore:      meta,
		chunkStore:     chunk,
		endpointsCount: cnt,
	}
	go svr.updateEndpointsCount()
	return svr
}

func (svr *Server) PutObject(ctx context.Context, bucket bucket.Bucket, fp fullpath.FullPath, reader io.Reader) error {
	ctime := time.Unix(time.Now().Unix(), 0)

	size := uint64(0)
	chunks := entry.Chunks{}

	md5Hash := md5.New()
	reader = io.TeeReader(reader, md5Hash)

	buf := mem.Allocate(partSize)
	defer mem.Free(buf)
	var uploadErr error
	for {
		n, err := reader.Read(buf)
		size += uint64(n)
		if size == 0 && err == io.EOF {
			break
		}
		var uploadFunc func(ctx context.Context, data []byte) (*entry.Chunk, error)
		if n <= smallObjectSize || svr.endpointsCount <= 2 {
			uploadFunc = svr.replicationUpload
		} else {
			uploadFunc = svr.reedSolomonUpload
		}
		if n > 0 {
			chunk, err := uploadFunc(ctx, buf[:n])
			if err != nil {
				uploadErr = err
				break
			} else {
				chunk.Offset = size - uint64(n)
				chunks = append(chunks, chunk)
			}
		}
		if err == io.EOF {
			break
		}
	}
	if uploadErr != nil {
		for _, chunk := range chunks {
			if chunk != nil {
				svr.clearFailedObject(ctx, chunk.Frags)
			}
		}
		return uploadErr
	}

	if size == 0 {
		err := svr.metaStore.InsertObject(ctx,
			&entry.Entry{
				FullPath: fp,
				Bucket:   bucket,
				Mtime:    ctime,
				Ctime:    ctime,
				Mode:     os.ModeDir,
				FileSize: size,
			})
		if err != nil {
			return err
		}
		return nil
	}

	ext := filepath.Ext(string(fp))
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType, reader = httputil.FileMimeDetect(reader)
	}

	ent := &entry.Entry{
		FullPath: fp,
		Bucket:   bucket,
		Mtime:    ctime,
		Ctime:    ctime,
		Mode:     os.ModePerm,
		Mime:     mimeType,
		Md5:      md5Hash.Sum(nil),
		FileSize: size,
		Chunks:   chunks,
	}

	err := svr.metaStore.InsertObject(ctx, ent)
	if err != nil {
		for _, chunk := range chunks {
			if chunk != nil {
				svr.clearFailedObject(ctx, chunk.Frags)
			}
		}
		return err
	}

	return nil
}

func (svr *Server) replicationUpload(ctx context.Context, data []byte) (*entry.Chunk, error) {
	hosts, err := svr.getDifferentHosts(ctx, replicationNum)
	if err != nil {
		return nil, err
	}

	lce := util.NewLimitedConcurrentExecutor(4)
	frags := make([]*entry.Frag, len(hosts))
	errChan := make(chan error, 1)
	for i_, host_ := range hosts {
		i, host := i_, host_
		lce.Execute(func() {
			var reader io.Reader
			reader = bytes.NewReader(data)
			md5Hash := md5.New()
			reader = io.TeeReader(reader, md5Hash)
			fid, err := svr.chunkStore.PutChunk(ctx, reader, host)
			if err != nil {
				errChan <- err
				return
			}
			frags[i] = &entry.Frag{
				Size:        uint64(len(data)),
				Id:          int64(i) + 1,
				Md5:         md5Hash.Sum(nil),
				IsDataShard: true,
				Fid:         fid,
			}
		})
	}
	go func() {
		lce.Wait()
		close(errChan)
	}()

	var merr error
	for err := range errChan {
		merr = multierr.Append(merr, err)
	}

	if merr != nil {
		svr.clearFailedObject(ctx, frags)
		return nil, merr
	}

	return &entry.Chunk{
		Size:          uint64(len(data)),
		Md5:           frags[0].Md5,
		Frags:         frags,
		IsReplication: len(hosts) >= 2,
	}, nil
}

func (svr *Server) reedSolomonUpload(ctx context.Context, data []byte) (*entry.Chunk, error) {
	hosts, err := svr.getDifferentHosts(ctx, reedSolomonMaxShard)
	if err != nil {
		return nil, err
	}

	md5Hash := md5.New()
	_, _ = md5Hash.Write(data)
	md5Ret := md5Hash.Sum(nil)
	dataShards, parityShards := dataParityShards[len(hosts)][0], dataParityShards[len(hosts)][1]
	shards := dataShards + parityShards
	encoder, _ := reedsolomon.New(dataShards, parityShards)

	// split data
	// calculate number of bytes per data shard.
	dataLen := len(data)
	perShard := (len(data) + dataShards - 1) / dataShards

	if cap(data) > len(data) {
		data = data[:cap(data)]
	}

	// only allocate memory if necessary
	var padding []byte
	if len(data) < (shards * perShard) {
		// calculate maximum number of full shards in `data` slice
		fullShards := len(data) / perShard
		buf := mem.Allocate(perShard * (shards - fullShards))
		defer mem.Free(buf)
		padding = buf
		copy(padding, data[perShard*fullShards:])
		data = data[0 : perShard*fullShards]
	} else {
		for i := dataLen; i < dataLen+dataShards; i++ {
			data[i] = 0
		}
	}

	// split into equal-length shards.
	d := make([][]byte, shards)
	i := 0
	for ; i < len(d) && len(data) >= perShard; i++ {
		d[i] = data[:perShard:perShard]
		data = data[perShard:]
	}

	for j := 0; i+j < len(d); j++ {
		d[i+j] = padding[:perShard:perShard]
		padding = padding[perShard:]
	}

	err = encoder.Encode(d)
	if err != nil {
		return nil, errors.ErrFSServer.WithMessagef("reed solomon encode faild: %v\n", err)
	}

	lce := util.NewLimitedConcurrentExecutor(8)
	errChan := make(chan error, 1)
	remainSize := dataLen
	frags := make([]*entry.Frag, len(hosts))
	for i := 0; i < shards; i++ {
		size := perShard
		if i < dataShards {
			size = min(len(d[i]), remainSize)
			remainSize -= size
		}
		obj := d[i][:size]
		lce.Execute(func() {
			var reader io.Reader
			reader = bytes.NewReader(obj)
			md5Hash := md5.New()
			reader = io.TeeReader(reader, md5Hash)
			fid, err := svr.chunkStore.PutChunk(ctx, reader, hosts[i])
			if err != nil {
				errChan <- err
				return
			}
			frags[i] = &entry.Frag{
				Size:        uint64(size),
				Id:          int64(i) + 1,
				Md5:         md5Hash.Sum(nil),
				IsDataShard: dataShardsPlan[shards][i],
				Fid:         fid,
			}
		})
	}
	go func() {
		lce.Wait()
		close(errChan)
	}()

	var merr error
	for err := range errChan {
		merr = multierr.Append(merr, err)
	}

	if merr != nil {
		svr.clearFailedObject(ctx, frags)
		return nil, merr
	}

	return &entry.Chunk{
		Size:          uint64(dataLen),
		Md5:           md5Ret,
		IsReplication: false,
		Frags:         frags,
	}, nil
}

func (svr *Server) GetEntry(ctx context.Context, bucket bucket.Bucket, fp fullpath.FullPath) (*entry.Entry, error) {
	return svr.metaStore.GetObject(ctx, bucket, fp)
}

func (svr *Server) GetObject(ctx context.Context, ent *entry.Entry, writer io.Writer) error {
	if !ent.Chunks.Verify() {
		return errors.ErrChunkMisalignment.WithMessagef("damaged file: %s", string(ent.FullPath))
	}

	for _, chunk := range ent.Chunks {
		if chunk.IsReplication {
			idx := rand.Intn(len(chunk.Frags))
			reader, err := svr.chunkStore.GetChunk(ctx, chunk.Frags[idx].Fid)
			if err != nil {
				return err
			}
			_, err = io.Copy(writer, reader)
			if err != nil {
				_ = reader.Close()
				return errors.ErrFSServer.WarpErr(err)
			}
			_ = reader.Close()
		} else {
			for _, frag := range chunk.Frags {
				if frag.IsDataShard {
					reader, err := svr.chunkStore.GetChunk(ctx, frag.Fid)
					if err != nil {
						return err
					}
					_, err = io.Copy(writer, reader)
					if err != nil {
						_ = reader.Close()
						return errors.ErrFSServer.WarpErr(err)
					}
					_ = reader.Close()
				}
			}
		}
	}
	return nil
}

func (svr *Server) ListObjects(ctx context.Context, bkt bucket.Bucket, fp fullpath.FullPath) ([]*entry.Entry, error) {
	return svr.metaStore.ListObjects(ctx, bkt, fp, false)
}

func (svr *Server) DeleteObject(ctx context.Context, bkt bucket.Bucket, fp fullpath.FullPath) error {
	return svr.metaStore.DeleteObjects(ctx, bkt, fp)
}

func (svr *Server) getDifferentHosts(ctx context.Context, size int) ([]string, error) {
	hosts, err := svr.chunkStore.GetEndpoints(ctx)
	if err != nil {
		return nil, err
	}
	rand.Shuffle(len(hosts), func(i, j int) {
		hosts[i], hosts[j] = hosts[j], hosts[i]
	})
	return hosts[:min(size, len(hosts))], nil
}

func (svr *Server) updateEndpointsCount() {
	ticker := time.NewTicker(1 * time.Minute)
	ctx := context.Background()
	for range ticker.C {
		hosts, err := svr.chunkStore.GetEndpoints(ctx)
		if err != nil {
			svr.endpointsCount = 0
		} else {
			svr.endpointsCount = len(hosts)
		}
	}
}

func (svr *Server) clearFailedObject(ctx context.Context, frags entry.Frags) {
	for _, frag := range frags {
		if frag != nil && frag.Fid != "" {
			_ = svr.chunkStore.DeleteChunk(ctx, frag.Fid)
		}
	}
}
