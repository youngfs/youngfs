package seaweedfs

import (
	"context"
	"github.com/go-playground/assert/v2"
	"testing"
	"youngfs/errors"
	"youngfs/kv/redis"
	"youngfs/vars"
)

func TestSeaweedFS_assignObject(t *testing.T) {
	kvStore := redis.NewKvStore(vars.RedisSocket, vars.RedisPassword, vars.RedisDatabase)
	client := NewStorageEngine(vars.SeaweedFSMaster, kvStore)
	info, err := client.assignObject(context.Background(), 5*1024)
	assert.Equal(t, err, nil)
	assert.Equal(t, info.Url, info.PublicUrl)
	assert.Equal(t, info.Count, int64(1))
}

func TestSeaweedFS_parseFid(t *testing.T) {
	ctx := context.Background()

	kvStore := redis.NewKvStore(vars.RedisSocket, vars.RedisPassword, vars.RedisDatabase)
	client := NewStorageEngine(vars.SeaweedFSMaster, kvStore)
	volumeId, fid, err := client.parseFid(ctx, "3,3fd41bd1da80")
	assert.Equal(t, volumeId, uint64(3))
	assert.Equal(t, fid, "3fd41bd1da80")
	assert.Equal(t, err, nil)

	volumeId, fid, err = client.parseFid(ctx, "3,3fd41bd1da80,3")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, errors.Is(err, errors.ErrServer), true)

	volumeId, fid, err = client.parseFid(ctx, "3fd41bd1da80")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, errors.Is(err, errors.ErrServer), true)

	volumeId, fid, err = client.parseFid(ctx, "")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, errors.Is(err, errors.ErrServer), true)

	volumeId, fid, err = client.parseFid(ctx, "-3,3fd41bd1da80")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, errors.Is(err, errors.ErrServer), true)

	volumeId, fid, err = client.parseFid(ctx, "3fd41bd1da80,3")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, errors.Is(err, errors.ErrServer), true)

	volumeId, fid, err = client.parseFid(ctx, "3fd41bd1da80.3")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, errors.Is(err, errors.ErrServer), true)
}
