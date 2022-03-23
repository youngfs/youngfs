package redis_store

import (
	"context"
	"github.com/go-playground/assert/v2"
	"icesos/command/vars"
	"icesos/errors"
	"icesos/kv"
	"icesos/util"
	"math/rand"
	"testing"
	"time"
)

func TestRedisStore_Num(t *testing.T) {
	client := NewRedisStore(vars.RedisHostPost, vars.RedisPassword, vars.RedisDatabase)
	key := "test_redis_num"
	ctx := context.Background()

	rand.Seed(time.Now().UnixNano())
	cnt := int64(0)

	ret, err := client.ClrNum(ctx, key)
	assert.Equal(t, ret, false)
	assert.Equal(t, err, kv.KvNotFound)

	for i := 0; i < 1024; i++ {
		rd := rand.Intn(2)
		if rd == 1 {
			cnt++
			ret, err := client.Incr(ctx, key)
			assert.Equal(t, ret, cnt)
			assert.Equal(t, err, nil)
			ret, err = client.GetNum(ctx, key)
			assert.Equal(t, ret, cnt)
			assert.Equal(t, err, nil)
		} else {
			cnt--
			ret, err := client.Decr(ctx, key)
			assert.Equal(t, ret, cnt)
			assert.Equal(t, err, nil)
			ret, err = client.GetNum(ctx, key)
			assert.Equal(t, ret, cnt)
			assert.Equal(t, err, nil)
		}
	}

	ret, err = client.ClrNum(ctx, key)
	assert.Equal(t, ret, true)
	assert.Equal(t, err, nil)

	ret, err = client.ClrNum(ctx, key)
	assert.Equal(t, ret, false)
	assert.Equal(t, err, kv.KvNotFound)

	err = client.SetNum(ctx, key, 128)
	cnt = int64(128)

	for i := 0; i < 1024; i++ {
		rd := rand.Intn(2)
		if rd == 1 {
			cnt++
			ret, err := client.Incr(ctx, key)
			assert.Equal(t, ret, cnt)
			assert.Equal(t, err, nil)
			ret, err = client.GetNum(ctx, key)
			assert.Equal(t, ret, cnt)
			assert.Equal(t, err, nil)
		} else {
			cnt--
			ret, err := client.Decr(ctx, key)
			assert.Equal(t, ret, cnt)
			assert.Equal(t, err, nil)
			ret, err = client.GetNum(ctx, key)
			assert.Equal(t, ret, cnt)
			assert.Equal(t, err, nil)
		}
	}

	err = client.SetNum(ctx, key, -127)
	cnt = int64(-127)

	for i := 0; i < 1024; i++ {
		rd := rand.Intn(2)
		if rd == 1 {
			cnt++
			ret, err := client.Incr(ctx, key)
			assert.Equal(t, ret, cnt)
			assert.Equal(t, err, nil)
			ret, err = client.GetNum(ctx, key)
			assert.Equal(t, ret, cnt)
			assert.Equal(t, err, nil)
		} else {
			cnt--
			ret, err := client.Decr(ctx, key)
			assert.Equal(t, ret, cnt)
			assert.Equal(t, err, nil)
			ret, err = client.GetNum(ctx, key)
			assert.Equal(t, ret, cnt)
			assert.Equal(t, err, nil)
		}
	}

	ret, err = client.ClrNum(ctx, key)
	assert.Equal(t, ret, true)
	assert.Equal(t, err, nil)

	err = client.KvPut(ctx, key, util.RandByte(1024))
	assert.Equal(t, err, nil)

	ret2, err := client.Incr(ctx, key)
	assert.Equal(t, ret2, int64(0))
	assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrKvSever])

	ret2, err = client.Decr(ctx, key)
	assert.Equal(t, ret2, int64(0))
	assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrKvSever])

	ret2, err = client.GetNum(ctx, key)
	assert.Equal(t, ret2, int64(0))
	assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrKvSever])

	ret, err = client.ClrNum(ctx, key)
	assert.Equal(t, ret, false)
	assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrKvSever])

	ret, err = client.KvDelete(ctx, key)
	assert.Equal(t, ret, true)
	assert.Equal(t, err, nil)

	ret2, err = client.GetNum(ctx, key)
	assert.Equal(t, ret2, int64(0))
	assert.Equal(t, err, kv.KvNotFound)

	for i := 0; i < 1024; i++ {
		cnt := rand.Int63()
		err := client.SetNum(ctx, key, cnt)
		assert.Equal(t, err, nil)
		ret, err := client.GetNum(ctx, key)
		assert.Equal(t, ret, cnt)
		assert.Equal(t, err, nil)
	}

	err = client.SetNum(ctx, key, rand.Int63())
	assert.Equal(t, err, nil)

	ret, err = client.ClrNum(ctx, key)
	assert.Equal(t, ret, true)
	assert.Equal(t, err, nil)

	ret2, err = client.GetNum(ctx, key)
	assert.Equal(t, ret2, int64(0))
	assert.Equal(t, err, kv.KvNotFound)

	ret, err = client.ClrNum(ctx, key)
	assert.Equal(t, ret, false)
	assert.Equal(t, err, kv.KvNotFound)

	ret2, err = client.GetNum(ctx, key)
	assert.Equal(t, ret2, int64(0))
	assert.Equal(t, err, kv.KvNotFound)

	ret3, err := client.KvGet(ctx, key)
	assert.Equal(t, ret3, nil)
	assert.Equal(t, err, kv.KvNotFound)

}
