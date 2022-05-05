package seaweedfs

import (
	"bytes"
	"context"
	"github.com/go-playground/assert/v2"
	"icesos/command/vars"
	"icesos/kv/redis"
	"icesos/util"
	"io/ioutil"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func TestSeaweedFS_PutObject(t *testing.T) {
	kvStore := redis.NewKvStore(vars.RedisHostPost, vars.RedisPassword, vars.RedisDatabase)
	client := NewStorageEngine(vars.MasterServer, kvStore)
	size := uint64(5 * 1024)
	ctx := context.Background()

	b := util.RandByte(size)

	fid, err := client.PutObject(ctx, size, bytes.NewReader(b), "")
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
	assert.Equal(t, httpBody, []byte{})

	hosts, err := client.GetHosts(ctx)
	fids := make([]string, 0)

	for i := 0; i < 16; i++ {
		b := util.RandByte(size)

		host := hosts[rand.Intn(len(hosts))]

		fid, err = client.PutObject(ctx, size, bytes.NewReader(b), host)
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

		httpBody, err = ioutil.ReadAll(resp.Body)
		assert.Equal(t, err, nil)
		assert.Equal(t, httpBody, b)

		fids = append(fids, fid)

		err = client.DeleteObject(ctx, fid)
		assert.Equal(t, err, nil)
	}

	time.Sleep(10 * time.Second)

	for _, fid := range fids {
		url, err = client.GetFidUrl(ctx, fid)
		assert.Equal(t, err, nil)

		resp, err := http.Get(url)
		assert.Equal(t, err, nil)
		defer func() {
			_ = resp.Body.Close()
		}()

		httpBody, err = ioutil.ReadAll(resp.Body)
		assert.Equal(t, err, nil)
		assert.Equal(t, httpBody, []byte{})
	}

}
