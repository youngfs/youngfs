package seaweedfs

import (
	"bytes"
	"context"
	"github.com/go-playground/assert/v2"
	"icesfs/command/vars"
	"icesfs/errors"
	"icesfs/kv/redis"
	"icesfs/log"
	"icesfs/util"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestSeaweedFS_DeleteObject(t *testing.T) {
	vars.UnitTest = true
	vars.Debug = true
	log.InitLogger()
	defer log.Sync()

	kvStore := redis.NewKvStore(vars.RedisSocket, vars.RedisPassword, vars.RedisDatabase)
	client := NewStorageEngine(vars.SeaweedFSMaster, kvStore)
	size := uint64(5 * 1024)
	ctx := context.Background()

	b := util.RandByte(size)

	fid, err := client.PutObject(ctx, size, bytes.NewReader(b), "", true)
	assert.Equal(t, err, nil)

	url, err := client.GetFidUrl(ctx, fid)
	assert.Equal(t, err, nil)

	resp1, err := http.Get(url)
	assert.Equal(t, err, nil)
	defer func() {
		_ = resp1.Body.Close()
	}()

	httpBody, err := ioutil.ReadAll(resp1.Body)
	assert.Equal(t, err, nil)
	assert.Equal(t, httpBody, b)

	req, err := http.NewRequest("DELETE", url, nil)
	assert.Equal(t, err, nil)

	resp, err := http.DefaultClient.Do(req)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusAccepted)

	resp2, err := http.Get(url)
	assert.Equal(t, err, nil)
	defer func() {
		_ = resp2.Body.Close()
	}()

	httpBody, err = ioutil.ReadAll(resp2.Body)
	assert.Equal(t, err, nil)
	assert.Equal(t, httpBody, []byte{})

	url, err = client.GetFidUrl(ctx, fid)
	assert.Equal(t, err, errors.GetAPIErr(errors.ErrObjectNotExist))
	assert.Equal(t, url, "")

	err = client.DeleteObject(ctx, fid)
	assert.Equal(t, err, nil)

	time.Sleep(3 * time.Second)
}
