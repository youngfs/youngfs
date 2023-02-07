package vfs

import (
	"bytes"
	"context"
	"github.com/go-playground/assert/v2"
	"os"
	"testing"
	"time"
	"youngfs/errors"
	"youngfs/fs/ec/ec_calc"
	"youngfs/fs/ec/ec_server"
	"youngfs/fs/ec/ec_store"
	"youngfs/fs/entry"
	"youngfs/fs/full_path"
	fs_set "youngfs/fs/set"
	"youngfs/fs/storage_engine/seaweedfs"
	"youngfs/kv/redis"
	"youngfs/util"
	"youngfs/vars"
)

func TestVFS(t *testing.T) {
	kvStore := redis.NewKvStore(vars.RedisSocket, vars.RedisPassword, vars.RedisDatabase)
	storageEngine := seaweedfs.NewStorageEngine(vars.SeaweedFSMaster, kvStore)
	ecStore := ec_store.NewECStore(kvStore, storageEngine)
	ecCalc := ec_calc.NewECCalc(ecStore, storageEngine)
	ecServer := ec_server.NewECServer(ecStore, ecCalc)
	vfs := NewVFS(kvStore, storageEngine, ecServer)

	set := fs_set.Set("test_vfs")
	mime := "application/octet-stream"
	size := uint64(5 * 1024)
	insertFiles := []full_path.FullPath{"/aa/bb/cc/dd", "/aa/bb/dd", "/aa/ee", "/ff", "/ll"}
	insertDirs := []full_path.FullPath{"/gg", "/bb/hh", "/aa/bb/ii", "/aa/bb/ee/jj", "/kk"}
	Files := []full_path.FullPath{"/aa/bb/cc/dd", "/aa/bb/dd", "/aa/ee", "/ff", "/ll"}
	Dirs := []full_path.FullPath{"/aa", "/aa/bb", "/aa/bb/cc", "/gg", "/bb/hh", "/aa/bb/ii", "/aa/bb/ee", "/aa/bb/ee/jj", "/kk"}
	time1 := time.Unix(time.Now().Unix(), 0)
	ctx := context.Background()

	fidMap := make(map[full_path.FullPath]string)

	for _, fp := range insertFiles {
		fid := putObject(t, ctx, vfs, size)
		fidMap[fp] = fid

		err := vfs.InsertObject(
			ctx,
			&entry.Entry{
				FullPath: fp,
				Set:      set,
				Mtime:    time1,
				Ctime:    time1,
				Mode:     os.ModePerm,
				Mime:     mime,
				FileSize: size,
				Fid:      fid,
			}, true)
		assert.Equal(t, err, nil)
	}

	for _, fp := range insertDirs {
		err := vfs.InsertObject(ctx,
			&entry.Entry{
				FullPath: fp,
				Set:      set,
				Mtime:    time1,
				Ctime:    time1,
				Mode:     os.ModeDir,
			}, true)
		assert.Equal(t, err, nil)
	}

	for _, fp := range Files {
		ent, err := vfs.GetObject(ctx, set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time1,
			Ctime:    time1,
			Mode:     os.ModePerm,
			Mime:     mime,
			FileSize: size,
			Fid:      fidMap[fp],
		})

		ent, err = vfs.getEntry(ctx, set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time1,
			Ctime:    time1,
			Mode:     os.ModePerm,
			Mime:     mime,
			FileSize: size,
			Fid:      fidMap[fp],
		})

		entries, err := vfs.ListObjects(ctx, set, fp)
		assert.Equal(t, errors.Is(err, errors.ErrInvalidPath), true)
		assert.Equal(t, entries, []entry.ListEntry{})
	}

	for _, fp := range Dirs {
		ent, err := vfs.GetObject(ctx, set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time1,
			Ctime:    time1,
			Mode:     os.ModeDir,
		})

		ent, err = vfs.getEntry(ctx, set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time1,
			Ctime:    time1,
			Mode:     os.ModeDir,
		})
	}

	ent, err := vfs.GetObject(ctx, set, "/")
	assert.Equal(t, errors.Is(err, errors.ErrInvalidPath), true)
	assert.Equal(t, ent, nil)

	ent, err = vfs.getEntry(ctx, set, "/")
	assert.Equal(t, errors.IsKvNotFound(err), true)
	assert.Equal(t, ent, nil)

	err = vfs.DeleteObject(ctx, set, full_path.FullPath("/"), false, time1)
	assert.Equal(t, errors.Is(err, errors.ErrInvalidDelete), true)

	for _, fp := range Files {
		ent, err := vfs.GetObject(ctx, set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time1,
			Ctime:    time1,
			Mode:     os.ModePerm,
			Mime:     mime,
			FileSize: size,
			Fid:      fidMap[fp],
		})

		ent, err = vfs.getEntry(ctx, set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time1,
			Ctime:    time1,
			Mode:     os.ModePerm,
			Mime:     mime,
			FileSize: size,
			Fid:      fidMap[fp],
		})

		entries, err := vfs.ListObjects(ctx, set, fp)
		assert.Equal(t, errors.Is(err, errors.ErrInvalidPath), true)
		assert.Equal(t, entries, []entry.ListEntry{})
	}

	for _, fp := range Dirs {
		ent, err := vfs.GetObject(ctx, set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time1,
			Ctime:    time1,
			Mode:     os.ModeDir,
		})

		ent, err = vfs.getEntry(ctx, set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time1,
			Ctime:    time1,
			Mode:     os.ModeDir,
		})
	}

	ent, err = vfs.GetObject(ctx, set, "/")
	assert.Equal(t, errors.Is(err, errors.ErrInvalidPath), true)
	assert.Equal(t, ent, nil)

	ent, err = vfs.getEntry(ctx, set, "/")
	assert.Equal(t, errors.IsKvNotFound(err), true)
	assert.Equal(t, ent, nil)

	time.Sleep(time.Duration(2) * time.Second)
	time2 := time.Unix(time.Now().Unix(), 0)

	err = vfs.InsertObject(ctx,
		&entry.Entry{
			FullPath: "/bb/hh",
			Set:      set,
			Mtime:    time2,
			Ctime:    time2,
			Mode:     os.ModeDir,
		}, false)
	assert.Equal(t, errors.Is(err, errors.ErrInvalidPath), true)

	err = vfs.InsertObject(ctx,
		&entry.Entry{
			FullPath: "/aa",
			Set:      set,
			Mtime:    time2,
			Ctime:    time2,
			Mode:     os.ModeDir,
		}, false)
	assert.Equal(t, errors.Is(err, errors.ErrInvalidPath), true)

	err = vfs.DeleteObject(ctx, set, full_path.FullPath("/aa/bb"), true, time2)
	assert.Equal(t, err, nil)

	Files1 := []full_path.FullPath{"/aa/ee", "/ff", "/ll"}       //mtime: time1 ctime: time1
	Dirs1 := []full_path.FullPath{"/gg", "/bb", "/bb/hh", "/kk"} //mtime: time1 ctime: time1
	Dirs2 := []full_path.FullPath{"/aa"}                         //mtime: time2 ctime: time1
	exist := make(map[full_path.FullPath]bool)
	for _, dir := range Files1 {
		exist[dir] = true
	}
	for _, dir := range Dirs1 {
		exist[dir] = true
	}
	for _, dir := range Dirs2 {
		exist[dir] = true
	}

	for _, fp := range Files {
		if exist[fp] {
			continue
		}
		ent, err := vfs.GetObject(ctx, set, fp)
		assert.Equal(t, errors.Is(err, errors.ErrInvalidPath), true)
		assert.Equal(t, ent, nil)

		ent, err = vfs.getEntry(ctx, set, fp)
		assert.Equal(t, errors.IsKvNotFound(err), true)
		assert.Equal(t, ent, nil)

		entries, err := vfs.ListObjects(ctx, set, fp)
		assert.Equal(t, errors.Is(err, errors.ErrInvalidPath), true)
		assert.Equal(t, entries, []entry.ListEntry{})
	}

	for _, fp := range Dirs {
		if exist[fp] {
			continue
		}
		ent, err := vfs.GetObject(ctx, set, fp)
		assert.Equal(t, errors.Is(err, errors.ErrInvalidPath), true)
		assert.Equal(t, ent, nil)

		ent, err = vfs.getEntry(ctx, set, fp)
		assert.Equal(t, errors.IsKvNotFound(err), true)
		assert.Equal(t, ent, nil)

		entries, err := vfs.ListObjects(ctx, set, fp)
		assert.Equal(t, errors.Is(err, errors.ErrInvalidPath), true)
		assert.Equal(t, entries, []entry.ListEntry{})
	}

	for _, fp := range Files1 {
		ent, err := vfs.GetObject(ctx, set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time1,
			Ctime:    time1,
			Mode:     os.ModePerm,
			Mime:     mime,
			FileSize: size,
			Fid:      fidMap[fp],
		})

		ent, err = vfs.getEntry(ctx, set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time1,
			Ctime:    time1,
			Mode:     os.ModePerm,
			Mime:     mime,
			FileSize: size,
			Fid:      fidMap[fp],
		})

		entries, err := vfs.ListObjects(ctx, set, fp)
		assert.Equal(t, errors.Is(err, errors.ErrInvalidPath), true)
		assert.Equal(t, entries, []entry.ListEntry{})
	}

	for _, fp := range Dirs1 {
		ent, err := vfs.GetObject(ctx, set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time1,
			Ctime:    time1,
			Mode:     os.ModeDir,
		})

		ent, err = vfs.getEntry(ctx, set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time1,
			Ctime:    time1,
			Mode:     os.ModeDir,
		})
	}

	for _, fp := range Dirs2 {
		ent, err := vfs.GetObject(ctx, set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time2,
			Ctime:    time1,
			Mode:     os.ModeDir,
		})

		ent, err = vfs.getEntry(ctx, set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time2,
			Ctime:    time1,
			Mode:     os.ModeDir,
		})
	}

	ent, err = vfs.GetObject(ctx, set, "/")
	assert.Equal(t, errors.Is(err, errors.ErrInvalidPath), true)
	assert.Equal(t, ent, nil)

	ent, err = vfs.getEntry(ctx, set, "/")
	assert.Equal(t, errors.IsKvNotFound(err), true)
	assert.Equal(t, ent, nil)

	time.Sleep(time.Duration(2) * time.Second)
	time3 := time.Unix(time.Now().Unix(), 0)

	fid := putObject(t, ctx, vfs, size)
	fidMap["/aa/ee/ll/mm"] = fid

	err = vfs.InsertObject(ctx,
		&entry.Entry{
			FullPath: "/aa/ee/ll/mm",
			Set:      set,
			Mtime:    time3,
			Ctime:    time3,
			Mode:     os.ModePerm,
			Mime:     mime,
			FileSize: size,
			Fid:      fid,
		}, true)
	assert.Equal(t, err, nil)

	err = vfs.InsertObject(ctx,
		&entry.Entry{
			FullPath: "/bb/hh/nn/oo",
			Set:      set,
			Mtime:    time3,
			Ctime:    time3,
			Mode:     os.ModeDir,
		}, true)
	assert.Equal(t, err, nil)

	err = vfs.InsertObject(ctx,
		&entry.Entry{
			FullPath: "/gg/pp",
			Set:      set,
			Mtime:    time3,
			Ctime:    time3,
			Mode:     os.ModeDir,
		}, true)
	assert.Equal(t, err, nil)

	fid = putObject(t, ctx, vfs, size)
	fidMap["/kk"] = fid

	err = vfs.InsertObject(ctx,
		&entry.Entry{
			FullPath: "/kk",
			Set:      set,
			Mtime:    time3,
			Ctime:    time3,
			Mode:     os.ModePerm,
			Mime:     mime,
			FileSize: size,
			Fid:      fid,
		}, true)
	assert.Equal(t, err, nil)

	err = vfs.InsertObject(ctx,
		&entry.Entry{
			FullPath: "/ll/rr",
			Set:      set,
			Mtime:    time3,
			Ctime:    time3,
			Mode:     os.ModeDir,
		}, true)
	assert.Equal(t, err, nil)

	Files1 = []full_path.FullPath{"/ff"}                                                                         //mtime:time1 ctime: time1
	Files2 := []full_path.FullPath{"/aa/ee/ll/mm", "/kk"}                                                        //mtime:time3 ctime: time3
	Dirs1 = []full_path.FullPath{"/bb"}                                                                          //mtime:time1 ctime: time1
	Dirs2 = []full_path.FullPath{"/aa", "/gg", "/bb/hh"}                                                         //mtime:time3 ctime: time1
	Dirs3 := []full_path.FullPath{"/aa/ee", "/aa/ee/ll", "/bb/hh/nn", "/bb/hh/nn/oo", "/gg/pp", "/ll", "/ll/rr"} //mtime:time3 ctime: time3
	Files = []full_path.FullPath{"/ff", "/aa/ee/ll/mm", "/kk", "/aa/bb/cc/dd", "/aa/bb/dd", "/aa/ee", "/ll"}
	Dirs = []full_path.FullPath{"/bb", "/aa", "/gg", "/bb/hh", "/aa/ee", "/aa/ee/ll", "/bb/hh/nn", "/bb/hh/nn/oo", "/gg/pp", "/ll", "/ll/rr",
		"/aa/bb", "/aa/bb/cc", "/aa/bb/ii", "/aa/bb/ee", "/aa/bb/ee/jj", "/kk"}
	exist = make(map[full_path.FullPath]bool)
	for _, dir := range Files1 {
		exist[dir] = true
	}
	for _, dir := range Files2 {
		exist[dir] = true
	}
	for _, dir := range Dirs1 {
		exist[dir] = true
	}
	for _, dir := range Dirs2 {
		exist[dir] = true
	}
	for _, dir := range Dirs3 {
		exist[dir] = true
	}

	for _, fp := range Files {
		if exist[fp] {
			continue
		}
		ent, err := vfs.GetObject(ctx, set, fp)
		assert.Equal(t, errors.Is(err, errors.ErrInvalidPath), true)
		assert.Equal(t, ent, nil)

		ent, err = vfs.getEntry(ctx, set, fp)
		assert.Equal(t, errors.IsKvNotFound(err), true)
		assert.Equal(t, ent, nil)

		entries, err := vfs.ListObjects(ctx, set, fp)
		assert.Equal(t, errors.Is(err, errors.ErrInvalidPath), true)
		assert.Equal(t, entries, []entry.ListEntry{})
	}

	for _, fp := range Dirs {
		if exist[fp] {
			continue
		}
		ent, err := vfs.GetObject(ctx, set, fp)
		assert.Equal(t, errors.Is(err, errors.ErrInvalidPath), true)
		assert.Equal(t, ent, nil)

		ent, err = vfs.getEntry(ctx, set, fp)
		assert.Equal(t, errors.IsKvNotFound(err), true)
		assert.Equal(t, ent, nil)

		entries, err := vfs.ListObjects(ctx, set, fp)
		assert.Equal(t, errors.Is(err, errors.ErrInvalidPath), true)
		assert.Equal(t, entries, []entry.ListEntry{})
	}

	for _, fp := range Files1 {
		ent, err := vfs.GetObject(ctx, set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time1,
			Ctime:    time1,
			Mode:     os.ModePerm,
			Mime:     mime,
			FileSize: size,
			Fid:      fidMap[fp],
		})

		ent, err = vfs.getEntry(ctx, set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time1,
			Ctime:    time1,
			Mode:     os.ModePerm,
			Mime:     mime,
			FileSize: size,
			Fid:      fidMap[fp],
		})

		entries, err := vfs.ListObjects(ctx, set, fp)
		assert.Equal(t, errors.Is(err, errors.ErrInvalidPath), true)
		assert.Equal(t, entries, []entry.ListEntry{})
	}

	for _, fp := range Files2 {
		ent, err := vfs.GetObject(ctx, set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time3,
			Ctime:    time3,
			Mode:     os.ModePerm,
			Mime:     mime,
			FileSize: size,
			Fid:      fidMap[fp],
		})

		ent, err = vfs.getEntry(ctx, set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time3,
			Ctime:    time3,
			Mode:     os.ModePerm,
			Mime:     mime,
			FileSize: size,
			Fid:      fidMap[fp],
		})

		entries, err := vfs.ListObjects(ctx, set, fp)
		assert.Equal(t, errors.Is(err, errors.ErrInvalidPath), true)
		assert.Equal(t, entries, []entry.ListEntry{})
	}

	for _, fp := range Dirs1 {
		ent, err := vfs.GetObject(ctx, set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time1,
			Ctime:    time1,
			Mode:     os.ModeDir,
		})

		ent, err = vfs.getEntry(ctx, set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time1,
			Ctime:    time1,
			Mode:     os.ModeDir,
		})
	}

	for _, fp := range Dirs2 {
		ent, err := vfs.GetObject(ctx, set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time3,
			Ctime:    time1,
			Mode:     os.ModeDir,
		})

		ent, err = vfs.getEntry(ctx, set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time3,
			Ctime:    time1,
			Mode:     os.ModeDir,
		})
	}

	for _, fp := range Dirs3 {
		ent, err := vfs.GetObject(ctx, set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time3,
			Ctime:    time3,
			Mode:     os.ModeDir,
		})

		ent, err = vfs.getEntry(ctx, set, fp)
		assert.Equal(t, err, nil)
		assert.Equal(t, ent, &entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time3,
			Ctime:    time3,
			Mode:     os.ModeDir,
		})
	}

	ent, err = vfs.GetObject(ctx, set, "/")
	assert.Equal(t, errors.Is(err, errors.ErrInvalidPath), true)
	assert.Equal(t, ent, nil)

	ent, err = vfs.getEntry(ctx, set, "/")
	assert.Equal(t, errors.IsKvNotFound(err), true)
	assert.Equal(t, ent, nil)

	err = vfs.DeleteObject(ctx, set, full_path.FullPath("/"), true, time3)
	assert.Equal(t, err, nil)

	for _, fp := range Files {
		ent, err := vfs.GetObject(ctx, set, fp)
		assert.Equal(t, errors.Is(err, errors.ErrInvalidPath), true)
		assert.Equal(t, ent, nil)

		ent, err = vfs.getEntry(ctx, set, fp)
		assert.Equal(t, errors.IsKvNotFound(err), true)
		assert.Equal(t, ent, nil)
	}

	for _, fp := range Dirs {
		ent, err := vfs.GetObject(ctx, set, fp)
		assert.Equal(t, errors.Is(err, errors.ErrInvalidPath), true)
		assert.Equal(t, ent, nil)

		ent, err = vfs.getEntry(ctx, set, fp)
		assert.Equal(t, errors.IsKvNotFound(err), true)
		assert.Equal(t, ent, nil)
	}

	ent, err = vfs.GetObject(ctx, set, "/")
	assert.Equal(t, errors.Is(err, errors.ErrInvalidPath), true)
	assert.Equal(t, ent, nil)

	ent, err = vfs.getEntry(ctx, set, "/")
	assert.Equal(t, errors.IsKvNotFound(err), true)
	assert.Equal(t, ent, nil)

	err = vfs.DeleteObject(ctx, set, full_path.FullPath("/"), true, time3)
	assert.Equal(t, err, nil)

	err = vfs.DeleteObject(ctx, set, full_path.FullPath("/aa"), true, time3)
	assert.Equal(t, errors.Is(err, errors.ErrInvalidPath), true)

	entries, err := vfs.ListObjects(ctx, set, full_path.FullPath("/"))
	assert.Equal(t, err, nil)
	assert.Equal(t, entries, []entry.ListEntry{})

	entries, err = vfs.ListObjects(ctx, set, full_path.FullPath("/aa"))
	assert.Equal(t, errors.Is(err, errors.ErrInvalidPath), true)
	assert.Equal(t, entries, []entry.ListEntry{})

	time.Sleep(3 * time.Second)
}

func putObject(t *testing.T, ctx context.Context, vfs *VFS, size uint64) string {
	b := util.RandByte(size)

	fid, err := vfs.storageEngine.PutObject(ctx, size, bytes.NewReader(b), "", true)
	assert.Equal(t, err, nil)

	return fid
}
