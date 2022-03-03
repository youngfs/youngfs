package server

import (
	"context"
	"crypto/md5"
	"icesos/entry"
	"icesos/errors"
	"icesos/filer"
	"icesos/full_path"
	"icesos/set"
	"icesos/storage_engine"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"time"
)

type Server struct {
	FilerStore    filer.FilerStore
	StorageEngine *storage_engine.StorageEngine
}

var Svr *Server

func NewServer(filer filer.FilerStore, storageEngine *storage_engine.StorageEngine) *Server {
	return &Server{
		FilerStore:    filer,
		StorageEngine: storageEngine,
	}
}

func (svr Server) PutObject(ctx context.Context, set set.Set, fp full_path.FullPath, size uint64, file multipart.File) error {
	ctime := time.Unix(time.Now().Unix(), 0)

	if size == 0 {
		err := svr.FilerStore.InsertObject(ctx,
			&entry.Entry{
				FullPath: fp,
				Set:      set,
				Ctime:    ctime,
				Mode:     os.ModeDir,
				FileSize: size,
			}, true)
		if err != nil {
			return err
		}
		return nil
	}

	fid, err := svr.StorageEngine.PutObject(ctx, size, file)
	if err != nil {
		return err
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return errors.ErrorCodeResponse[errors.ErrServer]
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return errors.ErrorCodeResponse[errors.ErrServer]
	}

	md5Ret := md5.Sum(b)

	err = file.Close()
	if err != nil {
		return errors.ErrorCodeResponse[errors.ErrServer]
	}

	err = svr.FilerStore.InsertObject(ctx,
		&entry.Entry{
			FullPath: fp,
			Set:      set,
			Ctime:    ctime,
			Mode:     os.ModePerm,
			Md5:      md5Ret,
			FileSize: size,
			Fid:      fid,
		}, true)
	if err != nil {
		return err
	}

	return nil
}

func (svr Server) GetObject(ctx context.Context, set set.Set, fp full_path.FullPath) (*entry.Entry, error) {
	return svr.FilerStore.GetObject(ctx, set, fp)
}

func (svr Server) ListObejcts(ctx context.Context, set set.Set, fp full_path.FullPath) ([]entry.ListEntry, error) {
	return svr.FilerStore.ListObjects(ctx, set, fp)
}

func (svr Server) DeleteObject(ctx context.Context, set set.Set, fp full_path.FullPath, recursive bool) error {
	return svr.FilerStore.DeleteObject(ctx, set, fp, recursive)
}
