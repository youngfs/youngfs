package entry

import (
	"github.com/go-playground/assert/v2"
	"github.com/go-redis/redis/v8"
	"icesos/command/vars"
	"icesos/full_path"
	"icesos/iam"
	"icesos/kv"
	"icesos/util"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestEntry(t *testing.T) {
	kv.Client.Initialize(vars.RedisHostPost, vars.RedisPassword, vars.RedisDatabase)

	fp := full_path.FullPath("/aa/bb/cc")
	set := iam.Set("test")

	Ct := time.Unix(time.Now().Unix(), 0) // windows: precision to s
	time.Sleep(time.Duration(2) * time.Second)

	entry := &Entry{
		FullPath: fp,
		Set:      set,
		Mtime:    time.Unix(time.Now().Unix(), 0), // windows: precision to s
		Ctime:    Ct,
		Mode:     os.ModePerm,
		Mime:     "",
		Md5:      util.RandMd5(),
		FileSize: uint64(rand.Int63()),
		VolumeId: uint64(rand.Int63()),
		Fid:      strconv.Itoa(rand.Int()),
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
