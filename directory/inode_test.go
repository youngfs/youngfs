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
		err := InsertInode(
			&Inode{
				FullPath: fp,
				Set:      set,
				Mtime:    time1,
				Ctime:    time1,
				Mode:     os.ModePerm,
			}, true)
		assert.Equal(t, err, nil)

		err = entry.InsertEntry(
			&entry.Entry{
				FullPath: fp,
				Set:      set,
				Mtime:    time1,
				Ctime:    time1,
				Mode:     os.ModePerm,
			})
		assert.Equal(t, err, nil)
	}

	for _, fp := range inodeDirs {
		err := InsertInode(
			&Inode{
				FullPath: fp,
				Set:      set,
				Mtime:    time1,
				Ctime:    time1,
				Mode:     os.ModeDir,
			}, true)
		assert.Equal(t, err, nil)

		err = entry.InsertEntry(
			&entry.Entry{
				FullPath: fp,
				Set:      set,
				Mtime:    time1,
				Ctime:    time1,
				Mode:     os.ModeDir,
			})
		assert.Equal(t, err, nil)
	}

	for _, fp := range entryFiles {
		nowEntry, err := entry.GetEntry(set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, nowEntry, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time1,
			Ctime:    time1,
			Mode:     os.ModePerm,
		})
	}

	for _, fp := range entryDirs {
		nowEntry, err := entry.GetEntry(set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, nowEntry, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time1,
			Ctime:    time1,
			Mode:     os.ModeDir,
		})
	}

	err := DeleteInodeAndEntry(set, full_path.FullPath("/"), time1, false)
	assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrInvalidDelete])

	for _, fp := range entryFiles {
		nowEntry, err := entry.GetEntry(set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, nowEntry, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time1,
			Ctime:    time1,
			Mode:     os.ModePerm,
		})
	}

	for _, fp := range entryDirs {
		nowEntry, err := entry.GetEntry(set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, nowEntry, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time1,
			Ctime:    time1,
			Mode:     os.ModeDir,
		})
	}

	time.Sleep(time.Duration(2) * time.Second)
	time2 := time.Unix(time.Now().Unix(), 0) // windows: precision to s

	err = DeleteInodeAndEntry(set, full_path.FullPath("/aa/bb"), time2, true)
	assert.Equal(t, err, nil)

	entryFiles1 := []full_path.FullPath{}                //ctime: time2  mtime:time2
	entryFiles2 := []full_path.FullPath{"/aa/ee", "/ff"} //ctime: time1  mtime:time1
	entryDirs1 := []full_path.FullPath{"/", "/aa"}       //ctime: time1  mtime:time2
	entryDirs2 := []full_path.FullPath{"/gg", "/aa/hh"}  //ctime: time1  mtime:time1
	set1 := make(map[full_path.FullPath]bool)
	for _, dir := range entryFiles1 {
		set1[dir] = true
	}
	for _, dir := range entryFiles2 {
		set1[dir] = true
	}

	set2 := make(map[full_path.FullPath]bool)
	for _, dir := range entryDirs1 {
		set2[dir] = true
	}
	for _, dir := range entryDirs2 {
		set2[dir] = true
	}

	for _, fp := range entryFiles {
		if set1[fp] {
			continue
		}
		nowEntry, err := entry.GetEntry(set, fp)
		assert.Equal(t, err, redis.Nil)
		assert.Equal(t, nowEntry, nil)
	}

	for _, fp := range entryDirs {
		if set2[fp] {
			continue
		}
		nowEntry, err := entry.GetEntry(set, fp)
		assert.Equal(t, err, redis.Nil)
		assert.Equal(t, nowEntry, nil)
	}

	for _, fp := range entryFiles1 {
		nowEntry, err := entry.GetEntry(set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, nowEntry, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time2,
			Ctime:    time2,
			Mode:     os.ModePerm,
		})
	}

	for _, fp := range entryFiles2 {
		nowEntry, err := entry.GetEntry(set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, nowEntry, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time1,
			Ctime:    time1,
			Mode:     os.ModePerm,
		})
	}

	for _, fp := range entryDirs1 {
		nowEntry, err := entry.GetEntry(set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, nowEntry, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time2,
			Ctime:    time1,
			Mode:     os.ModeDir,
		})
	}

	for _, fp := range entryDirs2 {
		nowEntry, err := entry.GetEntry(set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, nowEntry, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time1,
			Ctime:    time1,
			Mode:     os.ModeDir,
		})
	}

	time.Sleep(time.Duration(2) * time.Second)
	time3 := time.Unix(time.Now().Unix(), 0) // windows: precision to s

	err = InsertInode(
		&Inode{
			FullPath: "/aa/ee/kk",
			Set:      set,
			Mtime:    time3,
			Ctime:    time3,
			Mode:     os.ModePerm,
		}, true)
	assert.Equal(t, err, nil)

	err = entry.InsertEntry(
		&entry.Entry{
			FullPath: "/aa/ee/kk",
			Set:      set,
			Mtime:    time3,
			Ctime:    time3,
			Mode:     os.ModePerm,
		})
	assert.Equal(t, err, nil)

	err = InsertInode(
		&Inode{
			FullPath: "/aa/ll",
			Set:      set,
			Mtime:    time3,
			Ctime:    time3,
			Mode:     os.ModeDir,
		}, true)
	assert.Equal(t, err, nil)

	err = entry.InsertEntry(
		&entry.Entry{
			FullPath: "/aa/ll",
			Set:      set,
			Mtime:    time3,
			Ctime:    time3,
			Mode:     os.ModeDir,
		})
	assert.Equal(t, err, nil)

	entryFiles1 = []full_path.FullPath{"/aa/ee/kk"}                                   //ctime: time3  mtime:time3
	entryFiles2 = []full_path.FullPath{"/ff"}                                         //ctime: time1  mtime:time1
	entryDirs1 = []full_path.FullPath{"/", "/aa"}                                     //ctime: time1  mtime:time3
	entryDirs2 = []full_path.FullPath{"/gg", "/aa/hh"}                                //ctime: time1  mtime:time1
	entryDirs3 := []full_path.FullPath{"/aa/ee", "/aa/ll"}                            //ctime: time3  mtime:time3
	entryFiles = []full_path.FullPath{"/aa/bb/cc/dd", "/aa/bb/dd", "/ff", "aa/ee/kk"} // delete /aa/ee
	entryDirs = []full_path.FullPath{"/", "/aa", "/aa/bb", "/aa/bb/cc", "/gg", "/aa/hh", "/aa/bb/ii", "/aa/bb/ee", "/aa/bb/ee/jj", "/aa/ee", "/aa/ll"}
	set1 = make(map[full_path.FullPath]bool)
	for _, dir := range entryFiles1 {
		set1[dir] = true
	}
	for _, dir := range entryFiles2 {
		set1[dir] = true
	}

	set2 = make(map[full_path.FullPath]bool)
	for _, dir := range entryDirs1 {
		set2[dir] = true
	}
	for _, dir := range entryDirs2 {
		set2[dir] = true
	}
	for _, dir := range entryDirs3 {
		set2[dir] = true
	}

	for _, fp := range entryFiles {
		if set1[fp] {
			continue
		}
		nowEntry, err := entry.GetEntry(set, fp)
		assert.Equal(t, err, redis.Nil)
		assert.Equal(t, nowEntry, nil)
	}

	for _, fp := range entryDirs {
		if set2[fp] {
			continue
		}
		nowEntry, err := entry.GetEntry(set, fp)
		assert.Equal(t, err, redis.Nil)
		assert.Equal(t, nowEntry, nil)
	}

	for _, fp := range entryFiles1 {
		nowEntry, err := entry.GetEntry(set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, nowEntry, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time3,
			Ctime:    time3,
			Mode:     os.ModePerm,
		})
	}

	for _, fp := range entryFiles2 {
		nowEntry, err := entry.GetEntry(set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, nowEntry, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time1,
			Ctime:    time1,
			Mode:     os.ModePerm,
		})
	}

	for _, fp := range entryDirs1 {
		nowEntry, err := entry.GetEntry(set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, nowEntry, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time3,
			Ctime:    time1,
			Mode:     os.ModeDir,
		})
	}

	for _, fp := range entryDirs2 {
		nowEntry, err := entry.GetEntry(set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, nowEntry, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time1,
			Ctime:    time1,
			Mode:     os.ModeDir,
		})
	}

	for _, fp := range entryDirs3 {
		nowEntry, err := entry.GetEntry(set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, nowEntry, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time3,
			Ctime:    time3,
			Mode:     os.ModeDir,
		})
	}

	err = DeleteInodeAndEntry(set, full_path.FullPath("/"), time3, true)
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

	err = DeleteInodeAndEntry(set, full_path.FullPath("/"), time3, true)
	assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrInvalidDelete])

	err = DeleteInodeAndEntry(set, full_path.FullPath("/aa"), time3, true)
	assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrInvalidDelete])
}
