package redis

import (
	"context"
	"github.com/go-playground/assert/v2"
	"icesfs/command/vars"
	"icesfs/kv"
	"math/rand"
	"testing"
)

func TestRedis_ZSet(t *testing.T) {
	client := NewKvStore(vars.RedisSocket, vars.RedisPassword, vars.RedisDatabase)
	key := "test_redis_zset"
	ctx := context.Background()

	bList := make([]string, 26)
	for i := 0; i < 26; i++ {
		bList[i] = string(rune('a' + i))
	}

	for i := 0; i < 26; i++ {
		ret, err := client.ZRem(ctx, key, bList[i])
		assert.Equal(t, ret, false)
		assert.Equal(t, err, nil)
	}

	cnt, err := client.ZCard(ctx, key)
	assert.Equal(t, cnt, int64(0))
	assert.Equal(t, err, nil)

	for i := 0; i < 26; i++ {
		ret, err := client.ZRem(ctx, key, bList[i])
		assert.Equal(t, ret, false)
		assert.Equal(t, err, nil)
	}

	for i := 0; i < 16; i++ {
		err := client.ZAdd(ctx, key, bList[i])
		assert.Equal(t, err, nil)
	}

	cnt, err = client.ZCard(ctx, key)
	assert.Equal(t, cnt, int64(16))
	assert.Equal(t, err, nil)

	for i := 0; i < 26; i++ {
		ret, err := client.ZIsMember(ctx, key, bList[i])
		assert.Equal(t, ret, i < 16)
		assert.Equal(t, err, nil)
	}

	for i := 0; i < 26; i++ {
		err := client.ZAdd(ctx, key, bList[i])
		assert.Equal(t, err, nil)
	}

	cnt, err = client.ZCard(ctx, key)
	assert.Equal(t, cnt, int64(26))
	assert.Equal(t, err, nil)

	for i := 0; i < 26; i++ {
		ret, err := client.ZIsMember(ctx, key, bList[i])
		assert.Equal(t, ret, true)
		assert.Equal(t, err, nil)
	}

	for i := 0; i < 1000; i++ {
		x, y := rand.Int()%26, rand.Int()%26
		if x > y {
			x, y = y, x
		}
		ret, err := client.ZRangeByLex(ctx, key, string(rune('a'+x)), string(rune('a'+y)))
		assert.Equal(t, ret, bList[x:y])
		if x == y {
			assert.Equal(t, err, kv.NotFound)
		} else {
			assert.Equal(t, err, nil)
		}
	}

	bList2, err := client.ZRangeByLex(ctx, key, "b", "b")
	assert.Equal(t, bList2, []string{})
	assert.Equal(t, err, kv.NotFound)

	bList2, err = client.ZRangeByLex(ctx, key, "", "")
	assert.Equal(t, bList2, bList)
	assert.Equal(t, err, nil)

	ret, err := client.ZRemRangeByLex(ctx, key, string('a'+8), string('a'+20))
	assert.Equal(t, ret, true)
	assert.Equal(t, err, nil)

	cnt, err = client.ZCard(ctx, key)
	assert.Equal(t, cnt, int64(26-12))
	assert.Equal(t, err, nil)

	ret, err = client.ZRemRangeByLex(ctx, key, string('a'+8), string('a'+20))
	assert.Equal(t, ret, false)
	assert.Equal(t, err, nil)

	cnt, err = client.ZCard(ctx, key)
	assert.Equal(t, cnt, int64(26-12))
	assert.Equal(t, err, nil)

	for i := 0; i < 26; i++ {
		ret, err := client.ZIsMember(ctx, key, bList[i])
		assert.Equal(t, ret, i < 8 || i >= 20)
		assert.Equal(t, err, nil)
	}

	ret, err = client.ZRemRangeByLex(ctx, key, "", "")
	assert.Equal(t, ret, true)
	assert.Equal(t, err, nil)

	cnt, err = client.ZCard(ctx, key)
	assert.Equal(t, cnt, int64(0))
	assert.Equal(t, err, nil)

	ret, err = client.ZRemRangeByLex(ctx, key, "", "")
	assert.Equal(t, ret, false)
	assert.Equal(t, err, nil)

	cnt, err = client.ZCard(ctx, key)
	assert.Equal(t, cnt, int64(0))
	assert.Equal(t, err, nil)
}
