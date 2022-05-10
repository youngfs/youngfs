package vfs

import (
	"context"
	"github.com/go-playground/assert/v2"
	"icesos/command/vars"
	"icesos/ec/ec_calc"
	"icesos/ec/ec_server"
	"icesos/ec/ec_store"
	"icesos/entry"
	"icesos/full_path"
	"icesos/kv"
	"icesos/kv/redis"
	"icesos/set"
	"icesos/storage_engine/seaweedfs"
	"icesos/util"
	"os"
	"testing"
	"time"
)

func TestEntry(t *testing.T) {
	kvStore := redis.NewKvStore(vars.RedisHostPost, vars.RedisPassword, vars.RedisDatabase)
	storageEngine := seaweedfs.NewStorageEngine(vars.MasterServer, kvStore)
	ecStore := ec_store.NewEC(kvStore, storageEngine)
	ecCalc := ec_calc.NewECCalc(ecStore, storageEngine)
	ecServer := ec_server.NewECServer(ecStore, ecCalc)
	vfs := NewVFS(kvStore, storageEngine, ecServer)

	fp := full_path.FullPath("/aa/bb/cc")
	setName := set.Set("test_vfs_entry")
	ctx := context.Background()

	size := uint64(5 * 1024)

	fid := putObject(t, ctx, vfs, size)

	ent := &entry.Entry{
		FullPath: fp,
		Set:      setName,
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

	ent2, err := vfs.getEntry(ctx, setName, fp)
	assert.Equal(t, err, nil)
	assert.Equal(t, ent2, ent)
	assert.Equal(t, ent2.IsFile(), true)
	assert.Equal(t, ent2.IsDirectory(), false)

	err = vfs.deleteEntry(ctx, setName, fp)
	assert.Equal(t, err, nil)

	entry3, err := vfs.getEntry(ctx, setName, fp)
	assert.Equal(t, err, kv.NotFound)
	assert.Equal(t, entry3, nil)

	err = vfs.deleteEntry(ctx, setName, fp)
	assert.Equal(t, err, nil)
}
