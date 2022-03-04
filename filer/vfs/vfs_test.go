package vfs

import (
	"bytes"
	"context"
	"github.com/go-playground/assert/v2"
	"icesos/command/vars"
	"icesos/entry"
	"icesos/errors"
	"icesos/full_path"
	"icesos/kv"
	"icesos/kv/redis_store"
	"icesos/set"
	"icesos/storage_engine"
	"icesos/util"
	"os"
	"testing"
	"time"
)

func TestVFS(t *testing.T) {
	kvStore := redis_store.NewRedisStore(vars.RedisHostPost, vars.RedisPassword, vars.RedisDatabase)
	storageEngine := storage_engine.NewStorageEngine(vars.MasterServer)
	vfs := NewVFS(kvStore, storageEngine)

	setName := set.Set("test")
	mime := "text/plain"
	size := uint64(5 * 1024)
	inodesFiles := []full_path.FullPath{"/aa/bb/cc/dd", "/aa/bb/dd", "/aa/ee", "/ff"}
	inodeDirs := []full_path.FullPath{"/gg", "/aa/hh", "/aa/bb/ii", "/aa/bb/ee/jj"}
	entryFiles := []full_path.FullPath{"/aa/bb/cc/dd", "/aa/bb/dd", "/aa/ee", "/ff"}
	entryDirs := []full_path.FullPath{"/", "/aa", "/aa/bb", "/aa/bb/cc", "/gg", "/aa/hh", "/aa/bb/ii", "/aa/bb/ee", "/aa/bb/ee/jj"}
	ct1 := time.Unix(time.Now().Unix(), 0) // windows: precision to s
	ctx := context.Background()

	fidMap := make(map[full_path.FullPath]string)

	for _, fp := range inodesFiles {
		fid := putObject(t, ctx, vfs, size)
		fidMap[fp] = fid

		err := vfs.InsertObject(
			ctx,
			&entry.Entry{
				FullPath: fp,
				Set:      setName,
				Ctime:    ct1,
				Mode:     os.ModePerm,
				Mime:     mime,
				FileSize: size,
				Fid:      fid,
			}, true)
		assert.Equal(t, err, nil)
	}

	for _, fp := range inodeDirs {
		err := vfs.InsertObject(ctx,
			&entry.Entry{
				FullPath: fp,
				Set:      setName,
				Ctime:    ct1,
				Mode:     os.ModeDir,
			}, true)
		assert.Equal(t, err, nil)
	}

	for _, fp := range entryFiles {
		ent, err := vfs.GetObject(ctx, setName, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      setName,
			Ctime:    ct1,
			Mode:     os.ModePerm,
			Mime:     mime,
			FileSize: size,
			Fid:      fidMap[fp],
		})

		ent, err = vfs.getEntry(ctx, setName, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      setName,
			Ctime:    ct1,
			Mode:     os.ModePerm,
			Mime:     mime,
			FileSize: size,
			Fid:      fidMap[fp],
		})
	}

	for _, fp := range entryDirs {
		ent, err := vfs.GetObject(ctx, setName, fp)
		assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrInvalidPath])
		assert.Equal(t, ent, nil)

		ent, err = vfs.getEntry(ctx, setName, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      setName,
			Ctime:    ct1,
			Mode:     os.ModeDir,
		})
	}

	err := vfs.DeleteObject(ctx, setName, full_path.FullPath("/"), false)
	assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrInvalidDelete])

	for _, fp := range entryFiles {
		ent, err := vfs.GetObject(ctx, setName, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      setName,
			Ctime:    ct1,
			Mode:     os.ModePerm,
			Mime:     mime,
			FileSize: size,
			Fid:      fidMap[fp],
		})

		ent, err = vfs.getEntry(ctx, setName, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      setName,
			Ctime:    ct1,
			Mode:     os.ModePerm,
			Mime:     mime,
			FileSize: size,
			Fid:      fidMap[fp],
		})
	}

	for _, fp := range entryDirs {
		ent, err := vfs.GetObject(ctx, setName, fp)
		assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrInvalidPath])
		assert.Equal(t, ent, nil)

		ent, err = vfs.getEntry(ctx, setName, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      setName,
			Ctime:    ct1,
			Mode:     os.ModeDir,
		})
	}

	err = vfs.InsertObject(ctx,
		&entry.Entry{
			FullPath: "/aa/hh",
			Set:      setName,
			Ctime:    ct1,
			Mode:     os.ModeDir,
		}, false)
	assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrInvalidPath])

	err = vfs.DeleteObject(ctx, setName, full_path.FullPath("/aa/bb"), true)
	assert.Equal(t, err, nil)

	entryFiles1 := []full_path.FullPath{"/aa/ee", "/ff"}            //ctime: time1
	entryDirs1 := []full_path.FullPath{"/", "/aa", "/gg", "/aa/hh"} //ctime: time1
	set1 := make(map[full_path.FullPath]bool)
	for _, dir := range entryFiles1 {
		set1[dir] = true
	}

	set2 := make(map[full_path.FullPath]bool)
	for _, dir := range entryDirs1 {
		set2[dir] = true
	}

	for _, fp := range entryFiles {
		if set1[fp] {
			continue
		}
		ent, err := vfs.GetObject(ctx, setName, fp)
		assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrInvalidPath])
		assert.Equal(t, ent, nil)

		ent, err = vfs.getEntry(ctx, setName, fp)
		assert.Equal(t, err, kv.KvNotFound)
		assert.Equal(t, ent, nil)
	}

	for _, fp := range entryDirs {
		if set2[fp] {
			continue
		}
		ent, err := vfs.GetObject(ctx, setName, fp)
		assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrInvalidPath])
		assert.Equal(t, ent, nil)

		ent, err = vfs.getEntry(ctx, setName, fp)
		assert.Equal(t, err, kv.KvNotFound)
		assert.Equal(t, ent, nil)
	}

	for _, fp := range entryFiles1 {
		ent, err := vfs.GetObject(ctx, setName, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      setName,
			Ctime:    ct1,
			Mode:     os.ModePerm,
			Mime:     mime,
			FileSize: size,
			Fid:      fidMap[fp],
		})

		ent, err = vfs.getEntry(ctx, setName, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      setName,
			Ctime:    ct1,
			Mode:     os.ModePerm,
			Mime:     mime,
			FileSize: size,
			Fid:      fidMap[fp],
		})
	}

	for _, fp := range entryDirs1 {
		ent, err := vfs.GetObject(ctx, setName, fp)
		assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrInvalidPath])
		assert.Equal(t, ent, nil)

		ent, err = vfs.getEntry(ctx, setName, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      setName,
			Ctime:    ct1,
			Mode:     os.ModeDir,
		})
	}

	time.Sleep(time.Duration(2) * time.Second)
	ct2 := time.Unix(time.Now().Unix(), 0) // windows: precision to s

	fid := putObject(t, ctx, vfs, size)
	fidMap["/aa/ee/kk"] = fid

	err = vfs.InsertObject(ctx,
		&entry.Entry{
			FullPath: "/aa/ee/kk",
			Set:      setName,
			Ctime:    ct2,
			Mode:     os.ModePerm,
			Mime:     mime,
			FileSize: size,
			Fid:      fid,
		}, true)
	assert.Equal(t, err, nil)

	err = vfs.InsertObject(ctx,
		&entry.Entry{
			FullPath: "/aa/ll",
			Set:      setName,
			Ctime:    ct2,
			Mode:     os.ModeDir,
		}, true)
	assert.Equal(t, err, nil)

	err = vfs.InsertObject(ctx,
		&entry.Entry{
			FullPath: "/gg",
			Set:      setName,
			Ctime:    ct2,
			Mode:     os.ModeDir,
		}, true)
	assert.Equal(t, err, nil)

	entryFiles1 = []full_path.FullPath{"/ff"}                                         //ctime: time1
	entryFiles2 := []full_path.FullPath{"/aa/ee/kk"}                                  //ctime: time2
	entryDirs1 = []full_path.FullPath{"/", "/aa", "/aa/hh"}                           //ctime: time1
	entryDirs2 := []full_path.FullPath{"/gg", "/aa/ee", "/aa/ll"}                     //ctime: time2
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

	for _, fp := range entryFiles {
		if set1[fp] {
			continue
		}
		ent, err := vfs.GetObject(ctx, setName, fp)
		assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrInvalidPath])
		assert.Equal(t, ent, nil)

		ent, err = vfs.getEntry(ctx, setName, fp)
		assert.Equal(t, err, kv.KvNotFound)
		assert.Equal(t, ent, nil)
	}

	for _, fp := range entryDirs {
		if set2[fp] {
			continue
		}
		ent, err := vfs.GetObject(ctx, setName, fp)
		assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrInvalidPath])
		assert.Equal(t, ent, nil)

		ent, err = vfs.getEntry(ctx, setName, fp)
		assert.Equal(t, err, kv.KvNotFound)
		assert.Equal(t, ent, nil)
	}

	for _, fp := range entryFiles1 {
		ent, err := vfs.GetObject(ctx, setName, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      setName,
			Ctime:    ct1,
			Mode:     os.ModePerm,
			Mime:     mime,
			FileSize: size,
			Fid:      fidMap[fp],
		})

		ent, err = vfs.getEntry(ctx, setName, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      setName,
			Ctime:    ct1,
			Mode:     os.ModePerm,
			Mime:     mime,
			FileSize: size,
			Fid:      fidMap[fp],
		})
	}

	for _, fp := range entryFiles2 {
		ent, err := vfs.GetObject(ctx, setName, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      setName,
			Ctime:    ct2,
			Mode:     os.ModePerm,
			Mime:     mime,
			FileSize: size,
			Fid:      fidMap[fp],
		})

		ent, err = vfs.getEntry(ctx, setName, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      setName,
			Ctime:    ct2,
			Mode:     os.ModePerm,
			Mime:     mime,
			FileSize: size,
			Fid:      fidMap[fp],
		})
	}

	for _, fp := range entryDirs1 {
		ent, err := vfs.GetObject(ctx, setName, fp)
		assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrInvalidPath])
		assert.Equal(t, ent, nil)

		ent, err = vfs.getEntry(ctx, setName, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      setName,
			Ctime:    ct1,
			Mode:     os.ModeDir,
		})
	}

	for _, fp := range entryDirs2 {
		ent, err := vfs.GetObject(ctx, setName, fp)
		assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrInvalidPath])
		assert.Equal(t, ent, nil)

		ent, err = vfs.getEntry(ctx, setName, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      setName,
			Ctime:    ct2,
			Mode:     os.ModeDir,
		})
	}

	err = vfs.DeleteObject(ctx, setName, full_path.FullPath("/"), true)
	assert.Equal(t, err, nil)

	for _, fp := range entryFiles {
		ent, err := vfs.GetObject(ctx, setName, fp)
		assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrInvalidPath])
		assert.Equal(t, ent, nil)

		ent, err = vfs.getEntry(ctx, setName, fp)
		assert.Equal(t, err, kv.KvNotFound)
		assert.Equal(t, ent, nil)
	}

	for _, fp := range entryDirs {
		ent, err := vfs.GetObject(ctx, setName, fp)
		assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrInvalidPath])
		assert.Equal(t, ent, nil)

		ent, err = vfs.getEntry(ctx, setName, fp)
		assert.Equal(t, err, kv.KvNotFound)
		assert.Equal(t, ent, nil)
	}

	err = vfs.DeleteObject(ctx, setName, full_path.FullPath("/"), true)
	assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrInvalidPath])

	err = vfs.DeleteObject(ctx, setName, full_path.FullPath("/aa"), true)
	assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrInvalidPath])

	time.Sleep(3 * time.Second)
}

func putObject(t *testing.T, ctx context.Context, vfs *VFS, size uint64) string {
	b := util.RandByte(size)

	fid, err := vfs.storageEngine.PutObject(ctx, size, bytes.NewReader(b))
	assert.Equal(t, err, nil)

	return fid
}
