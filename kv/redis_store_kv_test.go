package kv

import (
	"github.com/go-playground/assert/v2"
	"github.com/go-redis/redis/v8"
	"icesos/command/vars"
	"icesos/util"
	"testing"
)

func TestRedisStore_Kv(t *testing.T) {
	Client.Initialize(vars.RedisHostPost, vars.RedisPassword, vars.RedisDatabase)
	key := "test_redis_kv"

	b, err := Client.KvGet(key)
	assert.Equal(t, err, redis.Nil)
	assert.Equal(t, b, nil)

	ret, err := Client.KvDelete(key)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, false)

	b = util.RandByte(1024)
	err = Client.KvPut(key, b)
	assert.Equal(t, err, nil)

	b2, err := Client.KvGet(key)
	assert.Equal(t, err, nil)
	assert.Equal(t, b2, b)

	b = util.RandByte(512)
	err = Client.KvPut(key, b)
	assert.Equal(t, err, nil)

	b2, err = Client.KvGet(key)
	assert.Equal(t, err, nil)
	assert.Equal(t, b2, b)

	ret, err = Client.KvDelete(key)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, true)

	b, err = Client.KvGet(key)
	assert.Equal(t, err, redis.Nil)
	assert.Equal(t, b, nil)

	ret, err = Client.KvDelete(key)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, false)
}
