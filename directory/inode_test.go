package directory

import (
	"github.com/go-playground/assert/v2"
	"github.com/go-redis/redis/v8"
	"icesos/command/vars"
	"icesos/entry"
	"icesos/errors"
	"icesos/full_path"
	"icesos/iam"
	"icesos/kv"
	"os"
	"testing"
	"time"
)

func TestInode(t *testing.T) {
	kv.Client.Initialize(vars.RedisHostPost, vars.RedisPassword, vars.RedisDatabase)

	set := iam.Set("test")
	inodesFiles := []full_path.FullPath{"/aa/bb/cc/dd", "/aa/bb/dd", "/aa/ee", "/ff"}
	inodeDirs := []full_path.FullPath{"/gg", "/aa/hh", "/aa/bb/ii", "/aa/bb/ee/jj"}
	entryFiles := []full_path.FullPath{"/aa/bb/cc/dd", "/aa/bb/dd", "/aa/ee", "/ff"}
	entryDirs := []full_path.FullPath{"/", "/aa", "/aa/bb", "/aa/bb/cc", "/gg", "/aa/hh", "/aa/bb/ii", "/aa/bb/ee", "/aa/bb/ee/jj"}
	time1 := time.Unix(time.Now().Unix(), 0) // windows: precision to s

	for _, fp := range inodesFiles {
		inode := &Inode{
			FullPath: fp,
			Set:      set,
			Time:     time1,
			Mode:     os.ModePerm,
		}
		err := InsertInode(inode, true)
		assert.Equal(t, err, nil)

		nowEntry := &entry.Entry{
			FullPath: fp,
			Set:      set,
			Time:     time1,
			Mode:     os.ModePerm,
		}
		err = entry.InsertEntry(nowEntry)
		assert.Equal(t, err, nil)
	}

	for _, fp := range inodeDirs {
		inode := &Inode{
			FullPath: fp,
			Set:      set,
			Time:     time.Unix(time.Now().Unix(), 0), // windows: precision to s
			Mode:     os.ModeDir,
		}
		err := InsertInode(inode, true)
		assert.Equal(t, err, nil)

		nowEntry := &entry.Entry{
			FullPath: fp,
			Set:      set,
			Time:     time1,
			Mode:     os.ModeDir,
		}
		err = entry.InsertEntry(nowEntry)
		assert.Equal(t, err, nil)
	}

	for _, fp := range entryFiles {
		nowEntry, err := entry.GetEntry(set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, nowEntry, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Time:     time1,
			Mode:     os.ModePerm,
		})
	}

	for _, fp := range entryDirs {
		nowEntry, err := entry.GetEntry(set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, nowEntry, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Time:     time1,
			Mode:     os.ModeDir,
		})
	}

	err := DeleteInodeAndEntry(set, full_path.FullPath("/"), false)
	assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrInvalidDelete])

	for _, fp := range entryFiles {
		nowEntry, err := entry.GetEntry(set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, nowEntry, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Time:     time1,
			Mode:     os.ModePerm,
		})
	}

	for _, fp := range entryDirs {
		nowEntry, err := entry.GetEntry(set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, nowEntry, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Time:     time1,
			Mode:     os.ModeDir,
		})
	}

	err = DeleteInodeAndEntry(set, full_path.FullPath("/"), true)
	assert.Equal(t, err, nil)

	for _, fp := range entryFiles {
		nowEntry, err := entry.GetEntry(set, fp)
		assert.Equal(t, err, redis.Nil)
		assert.Equal(t, nowEntry, nil)
	}

	for _, fp := range entryDirs {
		nowEntry, err := entry.GetEntry(set, fp)
		assert.Equal(t, err, redis.Nil)
		assert.Equal(t, nowEntry, nil)
	}
}
