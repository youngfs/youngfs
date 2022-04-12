package vfs

import (
	"context"
	"icesos/entry"
	"icesos/errors"
	"icesos/full_path"
	"icesos/kv"
	"icesos/set"
	"icesos/storage_engine"
	"time"
)

type VFS struct {
	kvStore       kv.KvStoreWithRedisMutex
	storageEngine storage_engine.StorageEngine
}

func NewVFS(kvStore kv.KvStoreWithRedisMutex, storageEngine storage_engine.StorageEngine) *VFS {
	return &VFS{
		kvStore:       kvStore,
		storageEngine: storageEngine,
	}
}

func (vfs *VFS) InsertObject(ctx context.Context, ent *entry.Entry, cover bool) error {
	if !ent.FullPath.IsLegalObjectName() {
		return errors.ErrorCodeResponse[errors.ErrIllegalObjectName]
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
		if err == kv.NotFound {
			return nil, errors.ErrorCodeResponse[errors.ErrInvalidPath]
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
			return errors.ErrorCodeResponse[errors.ErrInvalidDelete]
		}
	} else {
		ent, err := vfs.getEntry(ctx, set, fp)
		if err != nil {
			if err == kv.NotFound {
				return errors.ErrorCodeResponse[errors.ErrInvalidPath]
			}
			return err
		}

		inodeCnt, err := vfs.inodeCnt(ctx, set, fp)
		if err != nil {
			return err
		}

		if ent.IsDirectory() && recursive == false && inodeCnt != 0 {
			return errors.ErrorCodeResponse[errors.ErrInvalidDelete]
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
			if err == kv.NotFound {
				return []entry.ListEntry{}, errors.ErrorCodeResponse[errors.ErrInvalidPath]
			}
			return []entry.ListEntry{}, err
		}

		if ent.IsFile() {
			return []entry.ListEntry{}, errors.ErrorCodeResponse[errors.ErrInvalidPath]
		}
	}

	inodes, err := vfs.getInodeChs(ctx, set, fp)
	if err != nil {
		if err == kv.NotFound {
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
