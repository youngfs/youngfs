package server

import (
	"context"
	"crypto/md5"
	"io"
	"mime"
	"os"
	"path/filepath"
	"sort"
	"time"
	"youngfs/errors"
	"youngfs/fs/entry"
	"youngfs/fs/filer"
	"youngfs/fs/full_path"
	fs_set "youngfs/fs/set"
	"youngfs/fs/storage_engine"
	"youngfs/util"
)

type Server struct {
	filerStore    filer.FilerStore
	storageEngine storage_engine.StorageEngine
}

var svr *Server

func NewServer(filer filer.FilerStore, storageEngine storage_engine.StorageEngine) *Server {
	return &Server{
		filerStore:    filer,
		storageEngine: storageEngine,
	}
}

func PutObject(ctx context.Context, set fs_set.Set, fp full_path.FullPath, size uint64, reader io.Reader, compress bool) error {
	ctime := time.Unix(time.Now().Unix(), 0)

	if size == 0 {
		err := svr.filerStore.InsertObject(ctx,
			&entry.Entry{
				FullPath: fp,
				Set:      set,
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
		mimeType, reader = util.FileMimeDetect(reader)
	}

	ent := &entry.Entry{
		FullPath: fp,
		Set:      set,
		Mtime:    ctime,
		Ctime:    ctime,
		Mode:     os.ModePerm,
		Mime:     mimeType,
		FileSize: size,
	}
	chunks := entry.Chunks{}

	md5Hash := md5.New()
	reader = io.TeeReader(reader, md5Hash)

	for offset, id := uint64(0), int64(0); offset < size; offset += partSize {
		id++
		sz := util.Min(size-offset, partSize)
		md5Hash := md5.New()
		file := io.TeeReader(io.LimitReader(reader, int64(sz)), md5Hash)
		fid, err := svr.storageEngine.PutObject(ctx, sz, file, compress)
		if err != nil {
			for _, chunk := range chunks {
				for _, frag := range chunk.Frags {
					_ = svr.storageEngine.DeleteObject(ctx, frag.Fid)
				}
			}
			return err
		}
		chunks = append(chunks, entry.Chunk{
			Offset: offset,
			Size:   sz,
			Md5:    md5Hash.Sum(nil),
			Frags: entry.Frags{
				entry.Frag{
					Size:          sz,
					Id:            id,
					Md5:           md5Hash.Sum(nil),
					IsReplication: false,
					IsDataShard:   true,
					Fid:           fid,
				},
			},
		})
	}

	ent.Chunks = chunks
	ent.Md5 = md5Hash.Sum(nil)

	err := svr.filerStore.InsertObject(ctx, ent, true)
	if err != nil {
		for _, chunk := range chunks {
			for _, frag := range chunk.Frags {
				_ = svr.storageEngine.DeleteObject(ctx, frag.Fid)
			}
		}
		return err
	}

	return nil
}

func GetEntry(ctx context.Context, set fs_set.Set, fp full_path.FullPath) (*entry.Entry, error) {
	return svr.filerStore.GetObject(ctx, set, fp)
}

func GetObject(ctx context.Context, ent *entry.Entry, writer io.Writer) error {
	if !ent.Chunks.Verify() {
		return errors.ErrChunkMisalignment.WithMessagef("damaged file: %s", string(ent.FullPath))
	}
	for _, chunk := range ent.Chunks {
		sort.Sort(chunk.Frags)
		for _, frag := range chunk.Frags {
			if frag.IsDataShard {
				err := svr.storageEngine.GetObject(ctx, frag.Fid, writer)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func ListObjects(ctx context.Context, set fs_set.Set, fp full_path.FullPath) ([]entry.ListEntry, error) {
	return svr.filerStore.ListObjects(ctx, set, fp)
}

func DeleteObject(ctx context.Context, set fs_set.Set, fp full_path.FullPath, recursive bool) error {
	mtime := time.Unix(time.Now().Unix(), 0)
	return svr.filerStore.DeleteObject(ctx, set, fp, recursive, mtime)
}
