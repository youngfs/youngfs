package directory

import (
	"github.com/go-redis/redis/v8"
	"icesos/entry"
	"icesos/errors"
	"icesos/full_path"
	"icesos/iam"
	redis2 "icesos/kv/redis"
	"os"
	"time"
)

type Inode struct {
	full_path.FullPath             // full full_path
	iam.Set                        // own set_iam
	Mtime              time.Time   // time of last modification
	Ctime              time.Time   // time of creation
	Mode               os.FileMode // file mode
	Mime               string      // MIME type
	FileSize           uint64      // file size
}

func (inode *Inode) Key() string {
	return string(inode.Set) + string(inode.FullPath.Dir()) + inodeKv
}

func inodeKey(set iam.Set, fp full_path.FullPath) string {
	return string(set) + string(fp) + inodeKv
}

func (inode *Inode) IsDirectory() bool {
	return inode.Mode.IsDir()
}

func (inode *Inode) IsFile() bool {
	return inode.Mode.IsRegular()
}

func updateMtime(set iam.Set, fp full_path.FullPath, mtime time.Time) error {

	mutex := redis2.Client.NewMutex(string(set) + string(fp) + inodeLock)
	if err := mutex.Lock(); err != nil {
		return errors.ErrorCodeResponse[errors.ErrRedisSync]
	}
	defer func() {
		_, _ = mutex.Unlock()
	}()

	dirEntry, err := entry.GetEntry(set, fp)
	if err != nil {
		return err
	}

	oldMtime := dirEntry.Mtime //Record the Mtime before modification
	dirEntry.Mtime = mtime
	err = entry.InsertEntry(dirEntry)
	if err != nil {
		return err
	}

	if fp != inodeRoot {
		err := deleteInode(
			&Inode{
				FullPath: fp,
				Set:      set,
				Mtime:    oldMtime,
				Ctime:    dirEntry.Ctime,
				Mode:     os.ModeDir,
			})
		if err != nil {
			return err
		}

		err = insertInode(
			&Inode{
				FullPath: fp,
				Set:      set,
				Mtime:    mtime,
				Ctime:    dirEntry.Ctime,
				Mode:     os.ModeDir,
			})
		if err != nil {
			return err
		}
	}
	return nil
}

func insertInode(inode *Inode) error {
	b, err := inode.encodeProto()
	if err != nil {
		return err
	}

	err = redis2.Client.SAdd(inode.Key(), b)
	if err != nil {
		return err
	}

	return nil
}

func insertDir(inode *Inode, cover bool, dir full_path.FullPath) error {
	mutex := redis2.Client.NewMutex(string(inode.Set) + string(dir) + inodeLock)
	if err := mutex.Lock(); err != nil {
		return errors.ErrorCodeResponse[errors.ErrRedisSync]
	}
	defer func() {
		_, _ = mutex.Unlock()
	}()

	dirEntry, err := entry.GetEntry(inode.Set, dir)
	if err == redis.Nil {
		if dir != inode.FullPath {
			err = entry.InsertEntry(
				&entry.Entry{
					FullPath: dir,
					Set:      inode.Set,
					Mtime:    inode.Mtime,
					Ctime:    inode.Ctime,
					Mode:     os.ModeDir,
				},
			)
			if err != nil {
				return err
			}
		}

		if dir != inodeRoot {
			dirInode := &Inode{
				FullPath: dir,
				Set:      inode.Set,
				Mtime:    inode.Mtime,
				Ctime:    inode.Ctime,
				Mode:     os.ModeDir,
			}

			if dir == inode.FullPath {
				dirInode.Mode = inode.Mode
				dirInode.Mime = inode.Mime
				dirInode.FileSize = inode.FileSize
			}

			err := insertInode(dirInode)
			if err != nil {

				return err
			}
		}
	} else if err != nil {
		return err
	} else if dirEntry.IsFile() {
		if cover {
			//todo: if seaweedfs error, will not create inode
			err = deleteInodeAndEntry(inode.Set, dir, false)
			if err != nil {
				return err
			}

			if dir != inode.FullPath {
				err = entry.InsertEntry(
					&entry.Entry{
						FullPath: dir,
						Set:      inode.Set,
						Mtime:    inode.Mtime,
						Ctime:    inode.Ctime,
						Mode:     os.ModeDir,
					},
				)
				if err != nil {
					return err
				}
			}

			dirInode := &Inode{
				FullPath: dir,
				Set:      inode.Set,
				Mtime:    inode.Mtime,
				Ctime:    inode.Ctime,
				Mode:     os.ModeDir,
			}

			if dir == inode.FullPath {
				dirInode.Mode = inode.Mode
				dirInode.Mime = inode.Mime
				dirInode.FileSize = inode.FileSize
			}

			err = insertInode(dirInode)
			if err != nil {
				return err
			}
		} else {
			return errors.ErrorCodeResponse[errors.ErrInvalidPath]
		}
	}
	return nil
}

//todo: add insert inode and entry

func InsertInode(inode *Inode, cover bool) error {
	dirList := inode.SplitList()
	for _, dir := range dirList {
		err := insertDir(inode, cover, dir)
		if err != nil {
			return err
		}
	}

	if inode.FullPath != inodeRoot {
		err := updateMtime(inode.Set, inode.Dir(), inode.Mtime)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetInodes(set iam.Set, fp full_path.FullPath) ([]Inode, error) {
	key := inodeKey(set, fp)

	b, err := redis2.Client.SMembers(key)
	if err != nil {
		return nil, err
	}

	ret := make([]Inode, len(b))
	for i, v := range b {
		inode, err := decodeInodeProto(v)
		if err != nil {
			return nil, err
		}
		ret[i] = *inode
	}

	return ret, nil
}

func deleteInode(inode *Inode) error {
	b, err := inode.encodeProto()
	if err != nil {
		return err
	}

	_, err = redis2.Client.SRem(inode.Key(), b)
	if err != nil {
		return err
	}

	return nil
}

func deleteInodeAndEntry(set iam.Set, fp full_path.FullPath, lock bool) error {
	if lock {
		mutex := redis2.Client.NewMutex(string(set) + string(fp) + inodeLock)
		if err := mutex.Lock(); err != nil {
			return errors.ErrorCodeResponse[errors.ErrRedisSync]
		}
		defer func() {
			_, _ = mutex.Unlock()
		}()
	}

	inodes, err := GetInodes(set, fp)
	if err != nil {
		return err
	}

	for _, inode := range inodes {
		if inode.IsDirectory() {
			err = deleteInodeAndEntry(set, inode.FullPath, true)
			if err != nil {
				return err
			}
		} else {
			err = entry.DeleteEntry(set, inode.FullPath)
			if err != nil {
				return err
			}
		}
	}

	nowEntry, err := entry.GetEntry(set, fp)
	if err != nil {
		return err
	}

	if fp != inodeRoot {
		err := deleteInode(
			&Inode{
				FullPath: fp,
				Set:      set,
				Mtime:    nowEntry.Mtime,
				Ctime:    nowEntry.Ctime,
				Mode:     nowEntry.Mode,
				Mime:     nowEntry.Mime,
				FileSize: nowEntry.FileSize,
			})
		if err != nil {
			return err
		}
	}

	_, err = redis2.Client.SDelete(inodeKey(set, fp))
	if err != nil {
		return err
	}

	err = entry.DeleteEntry(set, fp)
	if err != nil {
		return err
	}

	return nil
}

func DeleteInodeAndEntry(set iam.Set, fp full_path.FullPath, mtime time.Time, recursive bool) error {

	inodeCnt, err := redis2.Client.SCard(inodeKey(set, fp))
	if err != nil {
		return err
	}

	_, err = entry.GetEntry(set, fp)
	if err != nil {
		if err == redis.Nil {
			return errors.ErrorCodeResponse[errors.ErrInvalidDelete]
		}
		return err
	}

	if recursive == false && inodeCnt != 0 {
		return errors.ErrorCodeResponse[errors.ErrInvalidDelete]
	}

	err = deleteInodeAndEntry(set, fp, true)
	if err != nil {
		return err
	}

	if fp != inodeRoot {
		err := updateMtime(set, fp.Dir(), mtime)
		if err != nil {
			return err
		}
	}

	return nil
}
