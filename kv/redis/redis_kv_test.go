package redis

import (
	"context"
	"github.com/go-playground/assert/v2"
	"icesos/command/vars"
	"icesos/kv"
	"icesos/util"
	"testing"
)

func TestRedis_Kv(t *testing.T) {
	client := NewKvStore(vars.RedisHostPost, vars.RedisPassword, vars.RedisDatabase)
	key := "test_redis_kv"
	ctx := context.Background()

	b, err := client.KvGet(ctx, key)
	assert.Equal(t, err, kv.NotFound)
	assert.Equal(t, b, nil)

	ret, err := client.KvDelete(ctx, key)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, false)

	b = util.RandByte(1024)
	err = client.KvPut(ctx, key, b)
	assert.Equal(t, err, nil)

	b2, err := client.KvGet(ctx, key)
	assert.Equal(t, err, nil)
	assert.Equal(t, b2, b)

	b = util.RandByte(512)
	err = client.KvPut(ctx, key, b)
	assert.Equal(t, err, nil)

	b2, err = client.KvGet(ctx, key)
	assert.Equal(t, err, nil)
	assert.Equal(t, b2, b)

	ret, err = client.KvDelete(ctx, key)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, true)

	b, err = client.KvGet(ctx, key)
	assert.Equal(t, err, kv.NotFound)
	assert.Equal(t, b, nil)

	ret, err = client.KvDelete(ctx, key)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, false)
}
