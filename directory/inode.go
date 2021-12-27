package directory

import (
	"github.com/go-redis/redis/v8"
	"icesos/entry"
	"icesos/errors"
	"icesos/full_path"
	"icesos/iam"
	"icesos/kv"
	"os"
	"time"
)

// Create / when server start

type Inode struct {
	full_path.FullPath             // full full_path
	iam.Set                        // own set_iam
	Time               time.Time   // time of creation
	Mode               os.FileMode // file mode
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

func InsertInode(inode *Inode, cover bool) error {
	dirList := inode.SplitList()
	for _, dir := range dirList {
		dirEntry, err := entry.GetEntry(inode.Set, dir)
		if err == redis.Nil {
			err = entry.InsertEntry(
				&entry.Entry{
					FullPath: dir,
					Set:      inode.Set,
					Time:     inode.Time,
					Mode:     os.ModeDir,
				},
			)
			if err != nil {
				return err
			}

			if dir != inodeRoot {
				dirInode := &Inode{
					FullPath: dir,
					Set:      inode.Set,
					Time:     inode.Time,
					Mode:     os.ModeDir,
				}

				if dir == inode.FullPath {
					dirInode.Mode = inode.Mode
				}

				b, err := dirInode.encodeProto()
				if err != nil {
					return err
				}

				err = kv.Client.SAdd(dirInode.Key(), b)
				if err != nil {
					return err
				}
			}

			continue
		} else if err != nil {
			return nil
		}

		if dirEntry.IsFile() {
			if cover {
				err = DeleteInodeAndEntry(inode.Set, dir, true)
				if err != nil {
					return err
				}

				dirInode := &Inode{
					FullPath: dir,
					Set:      inode.Set,
					Time:     inode.Time,
					Mode:     os.ModeDir,
				}

				if dir == inode.FullPath {
					dirInode.Mode = inode.Mode
				}

				b, err := dirInode.encodeProto()
				if err != nil {
					return err
				}

				err = kv.Client.SAdd(dirInode.Key(), b)
				if err != nil {
					return err
				}
			} else {
				return errors.ErrorCodeResponse[errors.ErrInvalidPath]
			}
		}
	}

	return nil
}

func GetInodes(set iam.Set, fp full_path.FullPath) ([]Inode, error) {
	key := inodeKey(set, fp)

	b, err := kv.Client.SMembers(key)
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

func DeleteInodeAndEntry(set iam.Set, fp full_path.FullPath, recursive bool) error {

	inodes, err := GetInodes(set, fp)
	if err != nil {
		return errors.ErrorCodeResponse[errors.ErrInvalidPath]
	}

	if recursive == false && len(inodes) != 0 {
		return errors.ErrorCodeResponse[errors.ErrInvalidDelete]
	}

	for _, inode := range inodes {
		if inode.IsDirectory() {
			err = DeleteInodeAndEntry(set, inode.FullPath, recursive)
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
		inode := &Inode{
			FullPath: fp,
			Set:      set,
			Time:     nowEntry.Time,
			Mode:     nowEntry.Mode,
		}

		b, err := inode.encodeProto()
		if err != nil {
			return err
		}

		_, err = kv.Client.SRem(inode.Key(), b)
		if err != nil {
			return err
		}
	}

	_, err = kv.Client.SDelete(inodeKey(set, fp))
	if err != nil {
		return err
	}

	err = entry.DeleteEntry(set, fp)
	if err != nil {
		return err
	}

	return nil
}
