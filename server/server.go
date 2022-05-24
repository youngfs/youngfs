package server

import (
	"context"
	"crypto/md5"
	"icesos/ec/ec_server"
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
	ecServer      *ec_server.ECServer
}

var svr *Server

func NewServer(filer filer.FilerStore, storageEngine storage_engine.StorageEngine, ecServer *ec_server.ECServer) *Server {
	return &Server{
		filerStore:    filer,
		storageEngine: storageEngine,
		ecServer:      ecServer,
	}
}

func PutObject(ctx context.Context, set set.Set, fp full_path.FullPath, size uint64, file io.Reader, compress bool) error {
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

	ent := &entry.Entry{
		FullPath: fp,
		Set:      set,
		Mtime:    ctime,
		Ctime:    ctime,
		Mode:     os.ModePerm,
		Mime:     mimeType,
		FileSize: size,
	}

	host, ecid, err := svr.ecServer.InsertObject(ctx, ent)
	if err != nil {
		return err
	}

	ent.ECid = ecid

	md5Hash := md5.New()
	file = io.TeeReader(file, md5Hash)

	fid := ""
	if host != "" {
		fid, err = svr.storageEngine.PutObject(ctx, size, file, fp.Name(), compress, host)
		if err != nil {
			_ = svr.ecServer.RecoverEC(ctx, ent)
			return err
		}
	} else {
		fid, err = svr.storageEngine.PutObject(ctx, size, file, fp.Name(), compress)
		if err != nil {
			_ = svr.ecServer.RecoverEC(ctx, ent)
			return err
		}
	}

	ent.Fid = fid
	ent.Md5 = md5Hash.Sum(nil)

	err = svr.filerStore.InsertObject(ctx, ent, true)
	if err != nil {
		_ = svr.storageEngine.DeleteObject(ctx, fid)
		_ = svr.ecServer.RecoverEC(ctx, ent)
		return err
	}

	err = svr.ecServer.ConfirmEC(ctx, ent)
	if err != nil {
		return err
	}

	err = svr.ecServer.ExecEC(ctx, ent.ECid)
	if err != nil {
		return err
	}

	return nil
}

func GetObject(ctx context.Context, set set.Set, fp full_path.FullPath) (*entry.Entry, error) {
	return svr.filerStore.GetObject(ctx, set, fp)
}

func ListObejcts(ctx context.Context, set set.Set, fp full_path.FullPath) ([]entry.ListEntry, error) {
	return svr.filerStore.ListObjects(ctx, set, fp)
}

func DeleteObject(ctx context.Context, set set.Set, fp full_path.FullPath, recursive bool) error {
	mtime := time.Unix(time.Now().Unix(), 0)
	return svr.filerStore.DeleteObject(ctx, set, fp, recursive, mtime)
}

func GetFidUrl(ctx context.Context, fid string) (string, error) {
	return svr.storageEngine.GetFidUrl(ctx, fid)
}

func RecoverObject(ctx context.Context, set set.Set, fp full_path.FullPath) error {
	ent, err := svr.filerStore.GetObject(ctx, set, fp)
	if err != nil {
		return err
	}

	frags, err := svr.ecServer.RecoverObject(ctx, ent)
	if err != nil {
		return err
	}
	return svr.filerStore.RecoverObject(ctx, frags)
}

func InsertSetRules(ctx context.Context, setRules *set.SetRules) error {
	return svr.ecServer.InsertSetRules(ctx, setRules)
}

func DeleteSetRules(ctx context.Context, set set.Set) error {
	return svr.ecServer.DeleteSetRules(ctx, set)
}

func GetSetRules(ctx context.Context, set set.Set) (*set.SetRules, error) {
	return svr.ecServer.GetSetRules(ctx, set)
}

func GetHosts(ctx context.Context) ([]string, error) {
	return svr.storageEngine.GetHosts(ctx)
}
