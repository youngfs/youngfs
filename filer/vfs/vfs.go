package vfs

import (
	"context"
	"icesos/entry"
	"icesos/errors"
	"icesos/full_path"
	"icesos/kv"
	"icesos/set"
	"icesos/storage_engine"
)

type VFS struct {
	kvStore       kv.KvStoreWithRedisMutex
	storageEngine *storage_engine.StorageEngine
}

func NewVFS(kvStore kv.KvStoreWithRedisMutex, storageEngine *storage_engine.StorageEngine) *VFS {
	return &VFS{
		kvStore:       kvStore,
		storageEngine: storageEngine,
	}
}

// after insert entry, insert inode
func (vfs VFS) InsertObject(ctx context.Context, ent *entry.Entry, cover bool) error {
	dirList := ent.SplitList()
	for _, dir := range dirList {
		err := vfs.insertInodeAndEntry(ctx, ent, dir, cover)
		if err != nil {
			return err
		}
	}
	return nil
}

func (vfs VFS) GetObject(ctx context.Context, set set.Set, fp full_path.FullPath) (*entry.Entry, error) {
	ent, err := vfs.getEntry(ctx, set, fp)
	if err != nil {
		if err == kv.KvNotFound {
			return nil, errors.ErrorCodeResponse[errors.ErrInvalidPath]
		}
		return nil, err
	}

	if ent.IsDirectory() {
		return nil, errors.ErrorCodeResponse[errors.ErrInvalidPath]
	}

	return ent, err
}

// after delete entry, delete inode
func (vfs VFS) DeleteObject(ctx context.Context, set set.Set, fp full_path.FullPath, recursive bool) error {
	_, err := vfs.getEntry(ctx, set, fp)
	if err != nil {
		if err == kv.KvNotFound {
			return errors.ErrorCodeResponse[errors.ErrInvalidDelete]
		}
		return err
	}

	inodeCnt, err := vfs.inodeCnt(ctx, set, fp)
	if err != nil {
		return err
	}

	if recursive == false && inodeCnt != 0 {
		return errors.ErrorCodeResponse[errors.ErrInvalidDelete]
	}

	err = vfs.deleteInodeAndEntry(ctx, set, fp, true)
	if err != nil {
		return err
	}
	return nil
}

func (vfs VFS) ListObjects(ctx context.Context, set set.Set, fp full_path.FullPath) ([]entry.ListEntry, error) {
	_, err := vfs.getEntry(ctx, set, fp)
	if err != nil {
		if err == kv.KvNotFound {
			return []entry.ListEntry{}, errors.ErrorCodeResponse[errors.ErrInvalidPath]
		}
		return []entry.ListEntry{}, err
	}

	inodes, err := vfs.getInodeChs(ctx, set, fp)
	if err != nil {
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

	return entry.ToListEntris(ret), nil
}
