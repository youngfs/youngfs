package seaweedfs

import (
	"bytes"
	"context"
	"github.com/go-playground/assert/v2"
	"io"
	"net/http"
	"testing"
	"time"
	"youngfs/errors"
	"youngfs/kv/redis"
	"youngfs/util"
	"youngfs/vars"
)

func TestSeaweedFS_DeleteObject(t *testing.T) {
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

	httpBody, err := io.ReadAll(resp1.Body)
	assert.Equal(t, err, nil)
	assert.Equal(t, httpBody, b)

	err = client.DeleteObject(ctx, fid)
	assert.Equal(t, err, nil)

	time.Sleep(3 * time.Second)

	url, err = client.GetFidUrl(ctx, fid)
	assert.Equal(t, errors.Is(err, errors.ErrObjectNotExist), true)
	assert.Equal(t, url, "")
}
