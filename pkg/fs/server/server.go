package server

import (
	"bytes"
	"context"
	"crypto/md5"
	"github.com/klauspost/reedsolomon"
	"github.com/youngfs/youngfs/pkg/errors"
	"github.com/youngfs/youngfs/pkg/fs/bucket"
	"github.com/youngfs/youngfs/pkg/fs/entry"
	"github.com/youngfs/youngfs/pkg/fs/filer"
	"github.com/youngfs/youngfs/pkg/fs/fullpath"
	"github.com/youngfs/youngfs/pkg/fs/storageengine"
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
	filerStore    filer.FilerStore
	storageEngine storageengine.StorageEngine
	hostCnt       int
}

var svr *Server

func NewServer(filer filer.FilerStore, storageEngine storageengine.StorageEngine) *Server {
	cnt := 0
	hosts, err := storageEngine.GetHosts(context.Background())
	if err == nil {
		cnt = len(hosts)
	}

	return &Server{
		filerStore:    filer,
		storageEngine: storageEngine,
		hostCnt:       cnt,
	}
}

func PutObject(ctx context.Context, bucket bucket.Bucket, fp fullpath.FullPath, reader io.Reader) error {
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
		if n <= smallObjectSize || svr.hostCnt <= 2 {
			uploadFunc = replicationUpload
		} else {
			uploadFunc = reedSolomonUpload
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
				clearFailedObject(ctx, chunk.Frags)
			}
		}
		return uploadErr
	}

	if size == 0 {
		err := svr.filerStore.InsertObject(ctx,
			&entry.Entry{
				FullPath: fp,
				Bucket:   bucket,
				Mtime:    ctime,
				Ctime:    ctime,
				Mode:     os.ModeDir,
				FileSize: size,
			}, true)
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

	err := svr.filerStore.InsertObject(ctx, ent, true)
	if err != nil {
		for _, chunk := range chunks {
			if chunk != nil {
				clearFailedObject(ctx, chunk.Frags)
			}
		}
		return err
	}

	return nil
}

func replicationUpload(ctx context.Context, data []byte) (*entry.Chunk, error) {
	hosts, err := getDifferentHosts(ctx, replicationNum)
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
			fid, err := svr.storageEngine.PutObject(ctx, uint64(len(data)), reader, host)
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
		clearFailedObject(ctx, frags)
		return nil, merr
	}

	return &entry.Chunk{
		Size:          uint64(len(data)),
		Md5:           frags[0].Md5,
		Frags:         frags,
		IsReplication: len(hosts) >= 2,
	}, nil
}

func reedSolomonUpload(ctx context.Context, data []byte) (*entry.Chunk, error) {
	hosts, err := getDifferentHosts(ctx, reedSolomonMaxShard)
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
		return nil, errors.ErrServer.WithMessagef("reed solomon encode faild: %v\n", err)
	}

	lce := util.NewLimitedConcurrentExecutor(8)
	errChan := make(chan error, 1)
	remainSize := dataLen
	frags := make([]*entry.Frag, len(hosts))
	for i_ := 0; i_ < shards; i_++ {
		i := i_
		size := perShard
		if i < dataShards {
			size = util.Min(len(d[i]), remainSize)
			remainSize -= size
		}
		obj := d[i][:size]
		lce.Execute(func() {
			var reader io.Reader
			reader = bytes.NewReader(obj)
			md5Hash := md5.New()
			reader = io.TeeReader(reader, md5Hash)
			fid, err := svr.storageEngine.PutObject(ctx, uint64(size), reader, hosts[i])
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
		clearFailedObject(ctx, frags)
		return nil, merr
	}

	return &entry.Chunk{
		Size:          uint64(dataLen),
		Md5:           md5Ret,
		IsReplication: false,
		Frags:         frags,
	}, nil
}

func GetEntry(ctx context.Context, bucket bucket.Bucket, fp fullpath.FullPath) (*entry.Entry, error) {
	return svr.filerStore.GetObject(ctx, bucket, fp)
}

func GetObject(ctx context.Context, ent *entry.Entry, writer io.Writer) error {
	if !ent.Chunks.Verify() {
		return errors.ErrChunkMisalignment.WithMessagef("damaged file: %s", string(ent.FullPath))
	}

	for _, chunk := range ent.Chunks {
		if chunk.IsReplication {
			idx := rand.Intn(len(chunk.Frags))
			err := svr.storageEngine.GetObject(ctx, chunk.Frags[idx].Fid, writer)
			if err != nil {
				return err
			}
		} else {
			for _, frag := range chunk.Frags {
				if frag.IsDataShard {
					err := svr.storageEngine.GetObject(ctx, frag.Fid, writer)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func ListObjects(ctx context.Context, bkt bucket.Bucket, fp fullpath.FullPath) ([]entry.ListEntry, error) {
	return svr.filerStore.ListObjects(ctx, bkt, fp)
}

func DeleteObject(ctx context.Context, bkt bucket.Bucket, fp fullpath.FullPath, recursive bool) error {
	mtime := time.Unix(time.Now().Unix(), 0)
	return svr.filerStore.DeleteObject(ctx, bkt, fp, recursive, mtime)
}

func getDifferentHosts(ctx context.Context, size int) ([]string, error) {
	hosts, err := svr.storageEngine.GetHosts(ctx)
	if err != nil {
		return nil, err
	}
	rand.Shuffle(len(hosts), func(i, j int) {
		hosts[i], hosts[j] = hosts[j], hosts[i]
	})
	return hosts[:util.Min(size, len(hosts))], nil
}

func (svr *Server) updateHostCnt() {
	ticker := time.NewTicker(5 * time.Second)
	ctx := context.Background()
	for range ticker.C {
		hosts, err := svr.storageEngine.GetHosts(ctx)
		if err != nil {
			svr.hostCnt = 0
		} else {
			svr.hostCnt = len(hosts)
		}
	}
}

func clearFailedObject(ctx context.Context, frags entry.Frags) {
	for _, frag := range frags {
		if frag != nil && frag.Fid != "" {
			_ = svr.storageEngine.DeleteObject(ctx, frag.Fid)
		}
	}
}
