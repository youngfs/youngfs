package vfs

import (
	"context"
	"icesos/command/vars"
	"icesos/entry"
	"icesos/errors"
	"icesos/full_path"
	"icesos/kv"
	"icesos/log"
	"icesos/set"
	"os"
	"time"
)

func inodeBelongKey(set set.Set, fp full_path.FullPath) string {
	return string(set) + string(fp.Dir()) + inodeKv
}

func inodeKey(set set.Set, fp full_path.FullPath) string {
	return string(set) + string(fp) + inodeKv
}

func inodeLockKey(set set.Set, fp full_path.FullPath) string {
	return string(set) + string(fp) + inodeLock
}

func (vfs *VFS) insertInodeFa(ctx context.Context, set set.Set, fp full_path.FullPath) error {
	if fp == inodeRoot {
		log.Errorw("insert inode: inode root", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), "set", set, "full path", fp)
		return errors.ErrorCodeResponse[errors.ErrServer]
	}

	err := vfs.kvStore.ZAdd(ctx, inodeBelongKey(set, fp), string(fp))
	if err != nil {
		return err
	}

	return nil
}

func (vfs *VFS) getInodeChs(ctx context.Context, set set.Set, fp full_path.FullPath) ([]full_path.FullPath, error) {
	fps, err := vfs.kvStore.ZRangeByLex(ctx, inodeKey(set, fp), "", "")
	if err != nil {
		return nil, err
	}

	ret := make([]full_path.FullPath, len(fps))
	for i, v := range fps {
		ret[i] = full_path.FullPath(v)
	}

	return ret, nil
}

func (vfs *VFS) deleteInodeFa(ctx context.Context, set set.Set, fp full_path.FullPath) error {
	if fp == inodeRoot {
		log.Errorw("delete inode: inode root", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), "set", set, "full path", fp)
		return errors.ErrorCodeResponse[errors.ErrServer]
	}

	_, err := vfs.kvStore.ZRem(ctx, inodeBelongKey(set, fp), string(fp))
	if err != nil {
		return err
	}

	return nil
}

func (vfs *VFS) deleteInodeChs(ctx context.Context, set set.Set, fp full_path.FullPath) error {
	_, err := vfs.kvStore.ZRemRangeByLex(ctx, inodeKey(set, fp), "", "")
	if err != nil {
		return err
	}
	return err
}

func (vfs *VFS) inodeCnt(ctx context.Context, set set.Set, fp full_path.FullPath) (int64, error) {
	return vfs.kvStore.ZCard(ctx, inodeKey(set, fp))
}

func (vfs *VFS) updateMtime(ctx context.Context, set set.Set, fp full_path.FullPath, mtime time.Time) error {
	mutex := vfs.kvStore.NewMutex(inodeLockKey(set, fp))
	if err := mutex.Lock(); err != nil {
		log.Errorw("update mtime lock error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "set", set, "full path", fp, "mtime", mtime)
		return errors.ErrorCodeResponse[errors.ErrRedisSync]
	}
	defer func() {
		_, _ = mutex.Unlock()
	}()

	ent, err := vfs.getEntry(ctx, set, fp)
	if err != nil {
		if err == kv.NotFound {
			log.Errorw("update mtime entry not found", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), "set", set, "full path", fp, "mtime", mtime)
			return errors.ErrorCodeResponse[errors.ErrServer]
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
func (vfs *VFS) insertInodeAndEntry(ctx context.Context, ent *entry.Entry, dir full_path.FullPath, cover bool) (bool, error) {
	mutex := vfs.kvStore.NewMutex(inodeLockKey(ent.Set, dir))
	if err := mutex.Lock(); err != nil {
		log.Errorw("insert inode and entry lock error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "entry", ent, "full path", dir, "cover", cover)
		return false, errors.ErrorCodeResponse[errors.ErrRedisSync]
	}
	defer func() {
		_, _ = mutex.Unlock()
	}()

	dirEnt, err := vfs.getEntry(ctx, ent.Set, dir)
	if err == kv.NotFound {
		if dir != ent.FullPath {
			err := vfs.insertEntry(ctx,
				&entry.Entry{
					FullPath: dir,
					Set:      ent.Set,
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

		err := vfs.insertInodeFa(ctx, ent.Set, dir)
		if err != nil {
			return false, err
		}

		return true, nil
	} else if err != nil {
		return false, err
	} else if dirEnt.IsFile() || dir == ent.FullPath {
		if cover {
			err = vfs.deleteInodeAndEntry(ctx, dirEnt.Set, dir, false)
			if err != nil {
				return false, err
			}

			if dir != ent.FullPath {
				err := vfs.insertEntry(ctx,
					&entry.Entry{
						FullPath: dir,
						Set:      ent.Set,
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

			err := vfs.insertInodeFa(ctx, ent.Set, dir)
			if err != nil {
				return false, err
			}

			return true, nil
		} else {
			return false, errors.ErrorCodeResponse[errors.ErrInvalidPath]
		}
	}
	return false, nil
}

func (vfs *VFS) deleteInodeAndEntry(ctx context.Context, set set.Set, fp full_path.FullPath, lock bool) error {
	if lock {
		mutex := vfs.kvStore.NewMutex(inodeLockKey(set, fp))
		if err := mutex.Lock(); err != nil {
			log.Errorw("delete inode and entry lock error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "set", set, "full path", fp, "lock", lock)
			return errors.ErrorCodeResponse[errors.ErrRedisSync]
		}
		defer func() {
			_, _ = mutex.Unlock()
		}()
	}
	inodes, err := vfs.getInodeChs(ctx, set, fp)
	if err != nil && err != kv.NotFound {
		return err
	}

	for _, inode := range inodes {
		ent, err := vfs.getEntry(ctx, set, inode)
		if err != nil {
			return err
		}

		if ent.IsDirectory() {
			err = vfs.deleteInodeAndEntry(ctx, set, ent.FullPath, true)
			if err != nil {
				return err
			}
		} else {
			err = vfs.deleteEntry(ctx, set, ent.FullPath)
			if err != nil {
				return err
			}
		}
	}

	// fp != /  delete fp.dir inode -> fp
	if fp != inodeRoot {
		err := vfs.deleteInodeFa(ctx, set, fp)
		if err != nil {
			return err
		}
	}

	err = vfs.deleteInodeChs(ctx, set, fp)
	if err != nil {
		return err
	}

	// fp != /  delete fp
	if fp != inodeRoot {
		err = vfs.deleteEntry(ctx, set, fp)
		if err != nil {
			return err
		}
	}

	return nil
}
