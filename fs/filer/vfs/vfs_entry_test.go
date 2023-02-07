package vfs

import (
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
	"youngfs/fs/id_generator/snow_flake"
	fs_set "youngfs/fs/set"
	"youngfs/fs/storage_engine/seaweedfs"
	"youngfs/kv/redis"
	"youngfs/util"
	"youngfs/vars"
)

func TestEntry(t *testing.T) {
	kvStore := redis.NewKvStore(vars.RedisSocket, vars.RedisPassword, vars.RedisDatabase)
	storageEngine := seaweedfs.NewStorageEngine(vars.SeaweedFSMaster, kvStore)
	ecStore := ec_store.NewECStore(kvStore, storageEngine, snow_flake.NewSnowFlake(0))
	ecCalc := ec_calc.NewECCalc(ecStore, storageEngine)
	ecServer := ec_server.NewECServer(ecStore, ecCalc)
	vfs := NewVFS(kvStore, storageEngine, ecServer)

	fp := full_path.FullPath("/aa/bb/cc")
	set := fs_set.Set("test_vfs_entry")
	ctx := context.Background()

	size := uint64(5 * 1024)

	fid := putObject(t, ctx, vfs, size)

	ent := &entry.Entry{
		FullPath: fp,
		Set:      set,
		Mtime:    time.Unix(time.Now().Unix(), 0),
		Ctime:    time.Unix(time.Now().Unix(), 0),
		Mode:     os.ModePerm,
		Mime:     "text/plain",
		Md5:      util.RandMd5(),
		FileSize: size,
		Fid:      fid,
	}

	assert.Equal(t, ent.IsFile(), true)
	assert.Equal(t, ent.IsDirectory(), false)

	err := vfs.insertEntry(ctx, ent)
	assert.Equal(t, err, nil)

	ent2, err := vfs.getEntry(ctx, set, fp)
	assert.Equal(t, err, nil)
	assert.Equal(t, ent2, ent)
	assert.Equal(t, ent2.IsFile(), true)
	assert.Equal(t, ent2.IsDirectory(), false)

	err = vfs.deleteEntry(ctx, set, fp)
	assert.Equal(t, err, nil)

	entry3, err := vfs.getEntry(ctx, set, fp)
	assert.Equal(t, errors.IsKvNotFound(err), true)
	assert.Equal(t, entry3, nil)

	err = vfs.deleteEntry(ctx, set, fp)
	assert.Equal(t, err, nil)
}
