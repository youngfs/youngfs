package server

import (
	"context"
	"crypto/md5"
	"icesos/entry"
	"icesos/filer"
	"icesos/full_path"
	"icesos/set"
	"icesos/storage_engine"
	"icesos/util"
	"io"
	"mime"
	"os"
	"path/filepath"
	"time"
)

type Server struct {
	filerStore    filer.FilerStore
	storageEngine storage_engine.StorageEngine
}

var Svr *Server

func NewServer(filer filer.FilerStore, storageEngine storage_engine.StorageEngine) *Server {
	return &Server{
		filerStore:    filer,
		storageEngine: storageEngine,
	}
}

func (svr Server) PutObject(ctx context.Context, set set.Set, fp full_path.FullPath, size uint64, file io.Reader) error {
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
		mimeType, file = util.FileMimeDetect(file)
	}

	md5Hash := md5.New()
	file = io.TeeReader(file, md5Hash)

	fid, err := svr.storageEngine.PutObject(ctx, size, file)
	if err != nil {
		return err
	}

	err = svr.filerStore.InsertObject(ctx,
		&entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    ctime,
			Ctime:    ctime,
			Mode:     os.ModePerm,
			Mime:     mimeType,
			Md5:      md5Hash.Sum(nil),
			FileSize: size,
			Fid:      fid,
		}, true)
	if err != nil {
		_ = svr.storageEngine.DeleteObject(ctx, fid)
		return err
	}

	return nil
}

func (svr Server) GetObject(ctx context.Context, set set.Set, fp full_path.FullPath) (*entry.Entry, error) {
	return svr.filerStore.GetObject(ctx, set, fp)
}

func (svr Server) ListObejcts(ctx context.Context, set set.Set, fp full_path.FullPath) ([]entry.ListEntry, error) {
	return svr.filerStore.ListObjects(ctx, set, fp)
}

func (svr Server) DeleteObject(ctx context.Context, set set.Set, fp full_path.FullPath, recursive bool) error {
	mtime := time.Unix(time.Now().Unix(), 0)
	return svr.filerStore.DeleteObject(ctx, set, fp, recursive, mtime)
}

func (svr Server) GetFidHost(ctx context.Context, fid string) (string, error) {
	return svr.storageEngine.GetFidHost(ctx, fid)
}
