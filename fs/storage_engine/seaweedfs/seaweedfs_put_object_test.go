package seaweedfs

import (
	"bytes"
	"context"
	"github.com/go-playground/assert/v2"
	"io"
	"math/rand"
	"net/http"
	"testing"
	"time"
	"youngfs/kv/redis"
	"youngfs/util"
	"youngfs/vars"
)

func TestSeaweedFS_PutObject(t *testing.T) {
	kvStore := redis.NewKvStore(vars.RedisSocket, vars.RedisPassword, vars.RedisDatabase)
	client := NewStorageEngine(vars.SeaweedFSMaster, kvStore)
	size := uint64(5 * 1024)
	ctx := context.Background()
	rand.Seed(time.Now().UnixNano())

	b := util.RandByte(size)

	fid, err := client.PutObject(ctx, size, bytes.NewReader(b), "", true, "")
	assert.Equal(t, err, nil)

	url, err := client.GetFidUrl(ctx, fid)
	assert.Equal(t, err, nil)

	resp1, err := http.Get(url)
	assert.Equal(t, err, nil)
	defer func() {
		_ = resp1.Body.Close()
	}()

	httpBody, err := io.ReadAll(resp1.Body)
	assert.Equal(t, err, nil)
	assert.Equal(t, httpBody, b)

	err = client.DeleteObject(ctx, fid)
	assert.Equal(t, err, nil)

	time.Sleep(3 * time.Second)

	resp2, err := http.Get(url)
	assert.Equal(t, err, nil)
	defer func() {
		_ = resp2.Body.Close()
	}()

	httpBody, err = io.ReadAll(resp2.Body)
	assert.Equal(t, err, nil)
	assert.Equal(t, httpBody, []byte{})

	hosts, err := client.GetHosts(ctx)
	urls := make([]string, 0)

	for i := 0; i < 128; i++ {
		b := util.RandByte(size)

		host := hosts[rand.Intn(len(hosts))]

		fid, err = client.PutObject(ctx, size, bytes.NewReader(b), "", rand.Intn(2) == rand.Intn(2), host)
		assert.Equal(t, err, nil)

		volumeId, _, err := client.parseFid(ctx, fid)
		assert.Equal(t, err, nil)

		host2, err := client.getVolumeHost(ctx, volumeId)
		assert.Equal(t, err, nil)
		assert.Equal(t, host2, host)

		url, err = client.GetFidUrl(ctx, fid)
		assert.Equal(t, err, nil)

		resp, err := http.Get(url)
		assert.Equal(t, err, nil)
		defer func() {
			_ = resp.Body.Close()
		}()

		httpBody, err = io.ReadAll(resp.Body)
		assert.Equal(t, err, nil)
		assert.Equal(t, httpBody, b)

		urls = append(urls, url)

		err = client.DeleteObject(ctx, fid)
		assert.Equal(t, err, nil)
	}

	time.Sleep(5 * time.Second)

	for _, url := range urls {
		resp, err := http.Get(url)
		assert.Equal(t, err, nil)
		defer func() {
			_ = resp.Body.Close()
		}()
		assert.Equal(t, resp.StatusCode, http.StatusNotFound)

		httpBody, err = io.ReadAll(resp.Body)
		assert.Equal(t, err, nil)
		assert.Equal(t, httpBody, []byte{})
	}

}
