package seaweedfs

import (
	"bytes"
	"context"
	"github.com/go-playground/assert/v2"
	"icesos/command/vars"
	"icesos/kv/redis"
	"icesos/util"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestSeaweedFS_AddLink(t *testing.T) {
	kvStore := redis.NewKvStore(vars.RedisSocket, vars.RedisPassword, vars.RedisDatabase)
	client := NewStorageEngine(vars.SeaweedFSMaster, kvStore)
	size := uint64(5 * 1024)
	ctx := context.Background()

	b := util.RandByte(size)

	fid, err := client.PutObject(ctx, size, bytes.NewReader(b), true)
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

	err = client.AddLink(ctx, fid)
	assert.Equal(t, err, nil)

	err = client.AddLink(ctx, fid)
	assert.Equal(t, err, nil)

	err = client.DeleteObject(ctx, fid)
	assert.Equal(t, err, nil)

	time.Sleep(3 * time.Second)

	resp2, err := http.Get(url)
	assert.Equal(t, err, nil)
	defer func() {
		_ = resp2.Body.Close()
	}()

	httpBody, err = ioutil.ReadAll(resp2.Body)
	assert.Equal(t, err, nil)
	assert.Equal(t, httpBody, b)

	err = client.DeleteObject(ctx, fid)
	assert.Equal(t, err, nil)

	time.Sleep(3 * time.Second)

	resp3, err := http.Get(url)
	assert.Equal(t, err, nil)
	defer func() {
		_ = resp3.Body.Close()
	}()

	httpBody, err = ioutil.ReadAll(resp3.Body)
	assert.Equal(t, err, nil)
	assert.Equal(t, httpBody, b)

	err = client.DeleteObject(ctx, fid)
	assert.Equal(t, err, nil)

	time.Sleep(3 * time.Second)

	resp4, err := http.Get(url)
	assert.Equal(t, err, nil)
	defer func() {
		_ = resp4.Body.Close()
	}()

	httpBody, err = ioutil.ReadAll(resp3.Body)
	assert.Equal(t, err, nil)
	assert.Equal(t, httpBody, []byte{})
}
