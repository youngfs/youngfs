package vfs

import (
	"context"
	"github.com/youngfs/youngfs/errors"
	"github.com/youngfs/youngfs/fs/bucket"
	"github.com/youngfs/youngfs/fs/entry"
	"github.com/youngfs/youngfs/fs/fullpath"
	"github.com/youngfs/youngfs/fs/storageengine"
	"github.com/youngfs/youngfs/kv"
	"time"
)

type VFS struct {
	kvStore       kv.KvSetStoreWithRedisMutex
	storageEngine storageengine.StorageEngine
}

func NewVFS(kvStore kv.KvSetStoreWithRedisMutex, storageEngine storageengine.StorageEngine) *VFS {
	return &VFS{
		kvStore:       kvStore,
		storageEngine: storageEngine,
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
				err := vfs.updateMtime(ctx, ent.Bucket, dir.Dir(), ent.Mtime)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (vfs *VFS) GetObject(ctx context.Context, bkt bucket.Bucket, fp fullpath.FullPath) (*entry.Entry, error) {
	ent, err := vfs.getEntry(ctx, bkt, fp)
	if err != nil {
		if errors.IsKvNotFound(err) {
			return nil, errors.ErrInvalidPath
		}
		return nil, err
	}

	return ent, err
}

// after delete entry, delete inode
func (vfs *VFS) DeleteObject(ctx context.Context, bkt bucket.Bucket, fp fullpath.FullPath, recursive bool, mtime time.Time) error {
	// if fp == / think fp is a folder
	if fp == inodeRoot {
		if recursive == false {
			return errors.ErrInvalidDelete
		}
	} else {
		ent, err := vfs.getEntry(ctx, bkt, fp)
		if err != nil {
			if errors.IsKvNotFound(err) {
				return errors.ErrInvalidPath
			}
			return err
		}

		inodeCnt, err := vfs.inodeCnt(ctx, bkt, fp)
		if err != nil {
			return err
		}

		if ent.IsDirectory() && recursive == false && inodeCnt != 0 {
			return errors.ErrInvalidDelete
		}
	}

	err := vfs.deleteInodeAndEntry(ctx, bkt, fp, true)
	if err != nil {
		return err
	}

	//include fp == / and fp.dir = /
	if fp.Dir() != inodeRoot {
		err := vfs.updateMtime(ctx, bkt, fp.Dir(), mtime)
		if err != nil {
			return err
		}
	}

	return nil
}

func (vfs *VFS) ListObjects(ctx context.Context, bkt bucket.Bucket, fp fullpath.FullPath) ([]entry.ListEntry, error) {
	//if fp != / check fp is dir
	if fp != inodeRoot {
		ent, err := vfs.getEntry(ctx, bkt, fp)
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

	inodes, err := vfs.getInodeChs(ctx, bkt, fp)
	if err != nil {
		if errors.IsKvNotFound(err) {
			return []entry.ListEntry{}, nil //not found return not err
		}
		return []entry.ListEntry{}, err
	}

	ret := make([]entry.Entry, len(inodes))
	for i, v := range inodes {
		ent, err := vfs.getEntry(ctx, bkt, v)
		if err != nil {
			return []entry.ListEntry{}, err
		}
		ret[i] = *ent
	}

	return entry.ToListEntries(ret), nil
}
