package entry

import (
	"bytes"
	"github.com/go-playground/assert/v2"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"icesos/command/vars"
	"icesos/full_path"
	"icesos/iam"
	"icesos/kv"
	"icesos/storage_engine"
	"icesos/util"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

func putObject(t *testing.T, size uint64) (uint64, string) {
	info, err := storage_engine.AssignObject(size)
	assert.Equal(t, err, nil)

	b := util.RandByte(size)
	req, err := http.NewRequest("PUT", "http://"+info.Url+"/"+info.Fid, bytes.NewReader(b))
	assert.Equal(t, err, nil)

	resp, err := http.DefaultClient.Do(req)
	assert.Equal(t, err, nil)

	httpBody, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, err, nil)

	putInfo := &storage_engine.PutObjectInfo{}
	err = jsoniter.Unmarshal(httpBody, putInfo)
	assert.Equal(t, err, nil)
	assert.Equal(t, putInfo.Size, size)

	resp, err = http.Get("http://" + info.Url + "/" + info.Fid)
	assert.Equal(t, err, nil)

	httpBody, err = ioutil.ReadAll(resp.Body)
	assert.Equal(t, err, nil)
	assert.Equal(t, httpBody, b)

	return storage_engine.SplitFid(info.Fid)
}

func TestEntry(t *testing.T) {
	kv.Client.Initialize(vars.RedisHostPost, vars.RedisPassword, vars.RedisDatabase)

	fp := full_path.FullPath("/aa/bb/cc")
	set := iam.Set("test")

	Ct := time.Unix(time.Now().Unix(), 0) // windows: precision to s
	time.Sleep(time.Duration(2) * time.Second)

	size := uint64(5 * 1024 * 1024)

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
