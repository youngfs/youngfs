package seaweedfs

import (
	"context"
	"github.com/go-playground/assert/v2"
	"icesos/command/vars"
	"icesos/errors"
	"icesos/kv/redis"
	"icesos/log"
	"testing"
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
	vars.UnitTest = true
	vars.Debug = true
	log.InitLogger()
	defer log.Sync()

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
	assert.Equal(t, err, errors.GetAPIErr(errors.ErrParseFid))

	volumeId, fid, err = client.parseFid(ctx, "3fd41bd1da80")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, err, errors.GetAPIErr(errors.ErrParseFid))

	volumeId, fid, err = client.parseFid(ctx, "")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, err, errors.GetAPIErr(errors.ErrParseFid))

	volumeId, fid, err = client.parseFid(ctx, "-3,3fd41bd1da80")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, err, errors.GetAPIErr(errors.ErrParseFid))

	volumeId, fid, err = client.parseFid(ctx, "3fd41bd1da80,3")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, err, errors.GetAPIErr(errors.ErrParseFid))

	volumeId, fid, err = client.parseFid(ctx, "3fd41bd1da80.3")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, err, errors.GetAPIErr(errors.ErrParseFid))
}
