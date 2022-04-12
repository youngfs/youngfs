package redis

import (
	"context"
	"github.com/go-playground/assert/v2"
	"icesos/command/vars"
	"icesos/kv"
	"icesos/util"
	"math/rand"
	"sort"
	"testing"
)

func TestRedis_Set(t *testing.T) {
	client := NewKvStore(vars.RedisHostPost, vars.RedisPassword, vars.RedisDatabase)
	key := "test_redis_set"
	ctx := context.Background()

	b := util.RandByte(512)

	bList := make([][]byte, 10)
	for i := 0; i < 10; i++ {
		bList[i] = util.RandByte(uint64(128 + rand.Int()%128))
	}

	// not add members
	bList2, err := client.SMembers(ctx, key)
	assert.Equal(t, err, kv.NotFound)
	assert.Equal(t, bList2, [][]byte{})

	cnt, err := client.SCard(ctx, key)
	assert.Equal(t, err, nil)
	assert.Equal(t, cnt, int64(0))

	ret, err := client.SRem(ctx, key, b)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, false)

	ret, err = client.SIsMember(ctx, key, b)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, false)

	ret, err = client.SDelete(ctx, key)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, false)

	//add members
	for _, m := range bList {
		err = client.SAdd(ctx, key, m)
		assert.Equal(t, err, nil)
	}

	bList2, err = client.SMembers(ctx, key)
	assert.Equal(t, err, nil)
	sort.Sort(util.BytesSlice(bList))
	sort.Sort(util.BytesSlice(bList2))
	assert.Equal(t, bList2, bList)

	cnt, err = client.SCard(ctx, key)
	assert.Equal(t, err, nil)
	assert.Equal(t, cnt, int64(10))

	ret, err = client.SRem(ctx, key, b)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, false)

	ret, err = client.SIsMember(ctx, key, b)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, false)

	for _, m := range bList {
		ret, err = client.SIsMember(ctx, key, m)
		assert.Equal(t, err, nil)
		assert.Equal(t, ret, true)
	}

	// delete some members
	for i := 7; i < 10; i++ {
		ret, err = client.SRem(ctx, key, bList[i])
		assert.Equal(t, err, nil)
		assert.Equal(t, ret, true)
	}

	bDel := bList[7:]
	bList = bList[:7]

	bList2, err = client.SMembers(ctx, key)
	assert.Equal(t, err, nil)
	sort.Sort(util.BytesSlice(bList))
	sort.Sort(util.BytesSlice(bList2))
	assert.Equal(t, bList2, bList)

	cnt, err = client.SCard(ctx, key)
	assert.Equal(t, err, nil)
	assert.Equal(t, cnt, int64(7))

	ret, err = client.SRem(ctx, key, b)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, false)

	for _, m := range bDel {
		ret, err = client.SRem(ctx, key, m)
		assert.Equal(t, err, nil)
		assert.Equal(t, ret, false)
	}

	ret, err = client.SIsMember(ctx, key, b)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, false)

	for _, m := range bList {
		ret, err = client.SIsMember(ctx, key, m)
		assert.Equal(t, err, nil)
		assert.Equal(t, ret, true)
	}

	for _, m := range bDel {
		ret, err = client.SIsMember(ctx, key, m)
		assert.Equal(t, err, nil)
		assert.Equal(t, ret, false)
	}

	// delete all members
	ret, err = client.SDelete(ctx, key)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, true)

	bList2, err = client.SMembers(ctx, key)
	assert.Equal(t, err, kv.NotFound)
	assert.Equal(t, bList2, [][]byte{})

	cnt, err = client.SCard(ctx, key)
	assert.Equal(t, err, nil)
	assert.Equal(t, cnt, int64(0))

	ret, err = client.SRem(ctx, key, b)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, false)

	for _, m := range bList {
		ret, err = client.SRem(ctx, key, m)
		assert.Equal(t, err, nil)
		assert.Equal(t, ret, false)
	}

	for _, m := range bDel {
		ret, err = client.SRem(ctx, key, m)
		assert.Equal(t, err, nil)
		assert.Equal(t, ret, false)
	}

	ret, err = client.SIsMember(ctx, key, b)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, false)

	for _, m := range bList {
		ret, err = client.SIsMember(ctx, key, m)
		assert.Equal(t, err, nil)
		assert.Equal(t, ret, false)
	}

	for _, m := range bDel {
		ret, err = client.SIsMember(ctx, key, m)
		assert.Equal(t, err, nil)
		assert.Equal(t, ret, false)
	}

	ret, err = client.SDelete(ctx, key)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, false)
}
