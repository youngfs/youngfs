package vfs

import (
	"context"
	"icesos/entry"
	"icesos/errors"
	"icesos/full_path"
	"icesos/kv"
	"icesos/set"
	"os"
)

func inodeKey(set set.Set, dir full_path.FullPath) string {
	return string(set) + string(dir.Dir()) + inodeKv
}

func (vfs VFS) insertInodeFa(ctx context.Context, set set.Set, fp full_path.FullPath) error {
	if fp == inodeRoot {
		return errors.ErrorCodeResponse[errors.ErrServer]
	}

	err := vfs.kvStore.ZAdd(ctx, inodeKey(set, fp), string(fp))
	if err != nil {
		return err
	}

	return nil
}

func (vfs VFS) getInodeChs(ctx context.Context, set set.Set, fp full_path.FullPath) ([]full_path.FullPath, error) {
	fps, err := vfs.kvStore.ZRangeByLex(ctx, string(set)+string(fp)+inodeKv, "", "")
	if err != nil {
		return nil, err
	}

	ret := make([]full_path.FullPath, len(fps))
	for i, v := range fps {
		ret[i] = full_path.FullPath(v)
	}

	return ret, nil
}

func (vfs VFS) deleteInodeFa(ctx context.Context, set set.Set, fp full_path.FullPath) error {
	if fp == inodeRoot {
		return errors.ErrorCodeResponse[errors.ErrServer]
	}

	_, err := vfs.kvStore.ZRem(ctx, inodeKey(set, fp), string(fp))
	if err != nil {
		return err
	}

	return nil
}

func (vfs VFS) deleteInodeChs(ctx context.Context, set set.Set, fp full_path.FullPath) error {
	_, err := vfs.kvStore.ZRemRangeByLex(ctx, string(set)+string(fp)+inodeKv, "", "")
	if err != nil {
		return err
	}
	return err
}

func (vfs VFS) inodeCnt(ctx context.Context, set set.Set, fp full_path.FullPath) (int64, error) {
	return vfs.kvStore.ZCard(ctx, inodeKey(set, fp))
}

func (vfs VFS) insertInodeAndEntry(ctx context.Context, ent *entry.Entry, dir full_path.FullPath, cover bool) error {
	mutex := vfs.kvStore.NewMutex(string(ent.Set) + string(dir) + inodeLock)
	if err := mutex.Lock(); err != nil {
		return errors.ErrorCodeResponse[errors.ErrRedisSync]
	}
	defer func() {
		_, _ = mutex.Unlock()
	}()

	dirEnt, err := vfs.getEntry(ctx, ent.Set, dir)
	if err == kv.KvNotFound {
		if dir != ent.FullPath {
			err := vfs.insertEntry(ctx,
				&entry.Entry{
					FullPath: dir,
					Set:      ent.Set,
					Ctime:    ent.Ctime,
					Mode:     os.ModeDir,
				},
			)
			if err != nil {
				return err
			}
		} else {
			err := vfs.insertEntry(ctx, ent)
			if err != nil {
				return err
			}
		}

		if dir != inodeRoot {
			err := vfs.insertInodeFa(ctx, ent.Set, dir)
			if err != nil {
				return err
			}
		}
	} else if err != nil {
		return err
	} else if dirEnt.IsFile() || dir == ent.FullPath {
		if cover {
			err = vfs.deleteInodeAndEntry(ctx, dirEnt.Set, dir, false)
			if err != nil {
				return err
			}

			if dir != ent.FullPath {
				err := vfs.insertEntry(ctx,
					&entry.Entry{
						FullPath: dir,
						Set:      ent.Set,
						Ctime:    ent.Ctime,
						Mode:     os.ModeDir,
					},
				)
				if err != nil {
					return err
				}
			} else {
				err := vfs.insertEntry(ctx, ent)
				if err != nil {
					return err
				}
			}

			if dir != inodeRoot {
				err := vfs.insertInodeFa(ctx, ent.Set, dir)
				if err != nil {
					return err
				}
			}
		} else {
			return errors.ErrorCodeResponse[errors.ErrInvalidPath]
		}
	}
	return nil
}

func (vfs VFS) deleteInodeAndEntry(ctx context.Context, set set.Set, fp full_path.FullPath, lock bool) error {
	if lock {
		mutex := vfs.kvStore.NewMutex(string(set) + string(fp) + inodeLock)
		if err := mutex.Lock(); err != nil {
			return errors.ErrorCodeResponse[errors.ErrRedisSync]
		}
		defer func() {
			_, _ = mutex.Unlock()
		}()
	}
	inodes, err := vfs.getInodeChs(ctx, set, fp)
	if err != nil && err != kv.KvNotFound {
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

	err = vfs.deleteEntry(ctx, set, fp)
	if err != nil {
		return err
	}

	return nil
}
