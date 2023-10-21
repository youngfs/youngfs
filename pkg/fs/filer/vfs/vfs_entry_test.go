package vfs

import (
	"context"
	"github.com/go-playground/assert/v2"
	"github.com/youngfs/youngfs/pkg/errors"
	"github.com/youngfs/youngfs/pkg/fs/bucket"
	"github.com/youngfs/youngfs/pkg/fs/entry"
	"github.com/youngfs/youngfs/pkg/fs/fullpath"
	"github.com/youngfs/youngfs/pkg/fs/storageengine/seaweedfs"
	"github.com/youngfs/youngfs/pkg/kv/redis"
	"github.com/youngfs/youngfs/pkg/util/randutil"
	"github.com/youngfs/youngfs/pkg/vars"
	"os"
	"testing"
	"time"
)

func TestEntry(t *testing.T) {
	kvStore := redis.NewKvStore(vars.RedisSocket, vars.RedisPassword, vars.RedisDatabase)
	storageEngine := seaweedfs.NewStorageEngine(vars.SeaweedFSMaster)
	vfs := NewVFS(kvStore, storageEngine)

	fp := fullpath.FullPath("/aa/bb/cc")
	bkt := bucket.Bucket("test")
	ctx := context.Background()

	size := uint64(20 * 1024)

	ent := &entry.Entry{
		FullPath: fp,
		Bucket:   bkt,
		Mtime:    time.Unix(time.Now().Unix(), 0),
		Ctime:    time.Unix(time.Now().Unix(), 0),
		Mode:     os.ModePerm,
		Mime:     "text/plain",
		Md5:      randutil.RandMd5(),
		FileSize: size,
		Chunks:   putObject(t, ctx, vfs, size),
	}

	assert.Equal(t, ent.IsFile(), true)
	assert.Equal(t, ent.IsDirectory(), false)

	err := vfs.insertEntry(ctx, ent)
	assert.Equal(t, err, nil)

	ent2, err := vfs.getEntry(ctx, bkt, fp)
	assert.Equal(t, err, nil)
	assert.Equal(t, ent2, ent)
	assert.Equal(t, ent2.IsFile(), true)
	assert.Equal(t, ent2.IsDirectory(), false)

	err = vfs.deleteEntry(ctx, bkt, fp)
	assert.Equal(t, err, nil)

	entry3, err := vfs.getEntry(ctx, bkt, fp)
	assert.Equal(t, errors.IsKvNotFound(err), true)
	assert.Equal(t, entry3, nil)

	err = vfs.deleteEntry(ctx, bkt, fp)
	assert.Equal(t, err, nil)

	time.Sleep(3 * time.Second)
}
