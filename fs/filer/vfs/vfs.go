package vfs

import (
	"context"
	"time"
	"youngfs/errors"
	"youngfs/fs/ec/ec_server"
	"youngfs/fs/ec/ec_store"
	"youngfs/fs/entry"
	"youngfs/fs/full_path"
	"youngfs/fs/set"
	"youngfs/fs/storage_engine"
	"youngfs/kv"
)

type VFS struct {
	kvStore       kv.KvSetStoreWithRedisMutex
	storageEngine storage_engine.StorageEngine
	ecServer      *ec_server.ECServer
}

func NewVFS(kvStore kv.KvSetStoreWithRedisMutex, storageEngine storage_engine.StorageEngine, ecServer *ec_server.ECServer) *VFS {
	return &VFS{
		kvStore:       kvStore,
		storageEngine: storageEngine,
		ecServer:      ecServer,
	}
}

func (vfs *VFS) InsertObject(ctx context.Context, ent *entry.Entry, cover bool) error {
	if !ent.FullPath.IsLegalObjectName() {
		return errors.ErrIllegalObjectName.Wrap("illegal object name: " + string(ent.FullPath))
	}

	dirList := ent.SplitList()
	isUpdateMtime := false
	for _, dir := range dirList[1:] {
		isCreate, err := vfs.insertInodeAndEntry(ctx, ent, dir, cover)
		if err != nil {
			return err
		}
		if !isUpdateMtime && isCreate {
			isUpdateMtime = true
			//only dir.dir == /
			if dir.Dir() != inodeRoot {
				err := vfs.updateMtime(ctx, ent.Set, dir.Dir(), ent.Mtime)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (vfs *VFS) GetObject(ctx context.Context, set set.Set, fp full_path.FullPath) (*entry.Entry, error) {
	ent, err := vfs.getEntry(ctx, set, fp)
	if err != nil {
		if errors.IsKvNotFound(err) {
			return nil, errors.ErrInvalidPath
		}
		return nil, err
	}

	return ent, err
}

// after delete entry, delete inode
func (vfs *VFS) DeleteObject(ctx context.Context, set set.Set, fp full_path.FullPath, recursive bool, mtime time.Time) error {
	// if fp == / think fp is a folder
	if fp == inodeRoot {
		if recursive == false {
			return errors.ErrInvalidDelete
		}
	} else {
		ent, err := vfs.getEntry(ctx, set, fp)
		if err != nil {
			if errors.IsKvNotFound(err) {
				return errors.ErrInvalidPath
			}
			return err
		}

		inodeCnt, err := vfs.inodeCnt(ctx, set, fp)
		if err != nil {
			return err
		}

		if ent.IsDirectory() && recursive == false && inodeCnt != 0 {
			return errors.ErrInvalidDelete
		}
	}

	err := vfs.deleteInodeAndEntry(ctx, set, fp, true)
	if err != nil {
		return err
	}

	//include fp == / and fp.dir = /
	if fp.Dir() != inodeRoot {
		err := vfs.updateMtime(ctx, set, fp.Dir(), mtime)
		if err != nil {
			return err
		}
	}

	return nil
}

func (vfs *VFS) ListObjects(ctx context.Context, set set.Set, fp full_path.FullPath) ([]entry.ListEntry, error) {
	//if fp != / check fp is dir
	if fp != inodeRoot {
		ent, err := vfs.getEntry(ctx, set, fp)
		if err != nil {
			if errors.IsKvNotFound(err) {
				return []entry.ListEntry{}, errors.ErrInvalidPath
			}
			return []entry.ListEntry{}, err
		}

		if ent.IsFile() {
			return []entry.ListEntry{}, errors.ErrInvalidPath
		}
	}

	inodes, err := vfs.getInodeChs(ctx, set, fp)
	if err != nil {
		if errors.IsKvNotFound(err) {
			return []entry.ListEntry{}, nil //not found return not err
		}
		return []entry.ListEntry{}, err
	}

	ret := make([]entry.Entry, len(inodes))
	for i, v := range inodes {
		ent, err := vfs.getEntry(ctx, set, v)
		if err != nil {
			return []entry.ListEntry{}, err
		}
		ret[i] = *ent
	}

	return entry.ToListEntries(ret), nil
}

func (vfs *VFS) RecoverObject(ctx context.Context, frags []ec_store.Frag) error {
	for _, frag := range frags {
		err := vfs.recoverEntry(ctx, &frag)
		if err != nil {
			return err
		}
	}

	return nil
}
