package vfs

import (
	"bytes"
	"github.com/go-playground/assert/v2"
	"github.com/go-redis/redis/v8"
	"icesos/command/vars"
	"icesos/full_path"
	"icesos/iam"
	redis2 "icesos/kv/redis"
	"icesos/storage_engine"
	"icesos/util"
	"os"
	"testing"
	"time"
)

func TestEntry(t *testing.T) {
	redis2.Client.Initialize(vars.RedisHostPost, vars.RedisPassword, vars.RedisDatabase)

	fp := full_path.FullPath("/aa/bb/cc")
	set := iam.Set("test")

	Ct := time.Unix(time.Now().Unix(), 0) // windows: precision to s
	time.Sleep(time.Duration(2) * time.Second)

	size := uint64(5 * 1024)

	volumeId, fid := putObject(t, size)

	entry := &Entry{
		FullPath: fp,
		Set:      set,
		Mtime:    time.Unix(time.Now().Unix(), 0), // windows: precision to s
		Ctime:    Ct,
		Mode:     os.ModePerm,
		Mime:     "",
		Md5:      util.RandMd5(),
		FileSize: size,
		VolumeId: volumeId,
		Fid:      fid,
	}

	assert.Equal(t, entry.IsFile(), true)
	assert.Equal(t, entry.IsDirectory(), false)

	err := InsertEntry(entry)
	assert.Equal(t, err, nil)

	entry2, err := GetEntry(set, fp)
	assert.Equal(t, err, nil)
	assert.Equal(t, entry2, entry)
	assert.Equal(t, entry2.IsFile(), true)
	assert.Equal(t, entry2.IsDirectory(), false)

	err = DeleteEntry(set, fp)
	assert.Equal(t, err, nil)

	entry3, err := GetEntry(set, fp)
	assert.Equal(t, err, redis.Nil)
	assert.Equal(t, entry3, nil)

	err = DeleteEntry(set, fp)
	assert.Equal(t, err, nil)
}

func putObject(t *testing.T, size uint64) (uint64, string) {
	b := util.RandByte(size)

	Fid, err := storage_engine.PutObject(size, bytes.NewReader(b))
	assert.Equal(t, err, nil)

	return storage_engine.SplitFid(Fid)
}
