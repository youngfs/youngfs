package vfs

import (
	"context"
	"github.com/youngfs/youngfs/pkg/errors"
	"github.com/youngfs/youngfs/pkg/fs/bucket"
	"github.com/youngfs/youngfs/pkg/fs/entry"
	"github.com/youngfs/youngfs/pkg/fs/fullpath"
	"os"
	"time"
)

func inodeBelongKey(bkt bucket.Bucket, fp fullpath.FullPath) string {
	return string(bkt) + string(fp.Dir()) + inodeKv
}

func inodeKey(bkt bucket.Bucket, fp fullpath.FullPath) string {
	return string(bkt) + string(fp) + inodeKv
}

func inodeLockKey(bkt bucket.Bucket, fp fullpath.FullPath) string {
	return string(bkt) + string(fp) + inodeLock
}

func (vfs *VFS) insertInodeFa(ctx context.Context, bkt bucket.Bucket, fp fullpath.FullPath) error {
	if fp == inodeRoot {
		return errors.ErrServer.Wrap("insert inode: inode root")
	}

	err := vfs.kvStore.ZAdd(ctx, inodeBelongKey(bkt, fp), string(fp))
	if err != nil {
		return err
	}

	return nil
}

func (vfs *VFS) getInodeChs(ctx context.Context, bkt bucket.Bucket, fp fullpath.FullPath) ([]fullpath.FullPath, error) {
	fps, err := vfs.kvStore.ZRangeByLex(ctx, inodeKey(bkt, fp), "", "")
	if err != nil {
		return nil, err
	}

	ret := make([]fullpath.FullPath, len(fps))
	for i, v := range fps {
		ret[i] = fullpath.FullPath(v)
	}

	return ret, nil
}

func (vfs *VFS) deleteInodeFa(ctx context.Context, bkt bucket.Bucket, fp fullpath.FullPath) error {
	if fp == inodeRoot {
		return errors.ErrServer.Wrap("delete inode: inode root")
	}

	_, err := vfs.kvStore.ZRem(ctx, inodeBelongKey(bkt, fp), string(fp))
	if err != nil {
		return err
	}

	return nil
}

func (vfs *VFS) deleteInodeChs(ctx context.Context, bkt bucket.Bucket, fp fullpath.FullPath) error {
	_, err := vfs.kvStore.ZRemRangeByLex(ctx, inodeKey(bkt, fp), "", "")
	if err != nil {
		return err
	}
	return err
}

func (vfs *VFS) inodeCnt(ctx context.Context, bkt bucket.Bucket, fp fullpath.FullPath) (int64, error) {
	return vfs.kvStore.ZCard(ctx, inodeKey(bkt, fp))
}

func (vfs *VFS) updateMtime(ctx context.Context, bkt bucket.Bucket, fp fullpath.FullPath, mtime time.Time) error {
	mutex := vfs.kvStore.NewMutex(inodeLockKey(bkt, fp))
	if err := mutex.Lock(); err != nil {
		return errors.ErrRedisSync.WithStack()
	}
	defer func() {
		_, _ = mutex.Unlock()
	}()

	ent, err := vfs.getEntry(ctx, bkt, fp)
	if err != nil {
		if errors.IsKvNotFound(err) {
			return errors.ErrServer.Wrap("update mtime: entry not found")
		}
		return err
	}

	ent.Mtime = mtime
	err = vfs.insertEntry(ctx, ent)
	if err != nil {
		return err
	}

	return nil
}

// return is create file or folder and err
func (vfs *VFS) insertInodeAndEntry(ctx context.Context, ent *entry.Entry, dir fullpath.FullPath, cover bool) (bool, error) {
	mutex := vfs.kvStore.NewMutex(inodeLockKey(ent.Bucket, dir))
	if err := mutex.Lock(); err != nil {
		return false, errors.ErrRedisSync.WithStack()
	}
	defer func() {
		_, _ = mutex.Unlock()
	}()

	dirEnt, err := vfs.getEntry(ctx, ent.Bucket, dir)
	if errors.IsKvNotFound(err) {
		if dir != ent.FullPath {
			err := vfs.insertEntry(ctx,
				&entry.Entry{
					FullPath: dir,
					Bucket:   ent.Bucket,
					Mtime:    ent.Mtime,
					Ctime:    ent.Ctime,
					Mode:     os.ModeDir,
				},
			)
			if err != nil {
				return false, err
			}
		} else {
			err := vfs.insertEntry(ctx, ent)
			if err != nil {
				return false, err
			}
		}

		err := vfs.insertInodeFa(ctx, ent.Bucket, dir)
		if err != nil {
			return false, err
		}

		return true, nil
	} else if err != nil {
		return false, err
	} else if dirEnt.IsFile() || dir == ent.FullPath {
		if cover {
			err = vfs.deleteInodeAndEntry(ctx, dirEnt.Bucket, dir, false)
			if err != nil {
				return false, err
			}

			if dir != ent.FullPath {
				err := vfs.insertEntry(ctx,
					&entry.Entry{
						FullPath: dir,
						Bucket:   ent.Bucket,
						Mtime:    ent.Mtime,
						Ctime:    ent.Ctime,
						Mode:     os.ModeDir,
					},
				)
				if err != nil {
					return false, err
				}
			} else {
				err := vfs.insertEntry(ctx, ent)
				if err != nil {
					return false, err
				}
			}

			err := vfs.insertInodeFa(ctx, ent.Bucket, dir)
			if err != nil {
				return false, err
			}

			return true, nil
		} else {
			return false, errors.ErrInvalidPath
		}
	}
	return false, nil
}

func (vfs *VFS) deleteInodeAndEntry(ctx context.Context, bkt bucket.Bucket, fp fullpath.FullPath, lock bool) error {
	if lock {
		mutex := vfs.kvStore.NewMutex(inodeLockKey(bkt, fp))
		if err := mutex.Lock(); err != nil {
			return errors.ErrRedisSync.WithStack()
		}
		defer func() {
			_, _ = mutex.Unlock()
		}()
	}
	inodes, err := vfs.getInodeChs(ctx, bkt, fp)
	if err != nil && !errors.IsKvNotFound(err) {
		return err
	}

	for _, inode := range inodes {
		ent, err := vfs.getEntry(ctx, bkt, inode)
		if err != nil {
			return err
		}

		if ent.IsDirectory() {
			err = vfs.deleteInodeAndEntry(ctx, bkt, ent.FullPath, true)
			if err != nil {
				return err
			}
		} else {
			err = vfs.deleteEntry(ctx, bkt, ent.FullPath)
			if err != nil {
				return err
			}
		}
	}

	// fp != /  delete fp.dir inode -> fp
	if fp != inodeRoot {
		err := vfs.deleteInodeFa(ctx, bkt, fp)
		if err != nil {
			return err
		}
	}

	err = vfs.deleteInodeChs(ctx, bkt, fp)
	if err != nil {
		return err
	}

	// fp != /  delete fp
	if fp != inodeRoot {
		err = vfs.deleteEntry(ctx, bkt, fp)
		if err != nil {
			return err
		}
	}

	return nil
}
