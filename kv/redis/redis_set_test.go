package redis

import (
	"bytes"
	"context"
	"github.com/go-playground/assert/v2"
	"math/rand"
	"sort"
	"testing"
	"youngfs/errors"
	"youngfs/util/randutil"
	"youngfs/vars"
)

type bytesSlice [][]byte

func (b bytesSlice) Len() int {
	return len(b)
}

func (b bytesSlice) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b bytesSlice) Less(i, j int) bool {
	return bytes.Compare(b[i], b[j]) < 0
}

func TestRedis_Set(t *testing.T) {
	client := NewKvStore(vars.RedisSocket, vars.RedisPassword, vars.RedisDatabase)
	key := "test_redis_set"
	ctx := context.Background()

	b := randutil.RandByte(512)

	bList := make([][]byte, 10)
	for i := 0; i < 10; i++ {
		bList[i] = randutil.RandByte(uint64(128 + rand.Int()%128))
	}

	// not add members
	bList2, err := client.SMembers(ctx, key)
	assert.Equal(t, errors.IsKvNotFound(err), true)
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
	sort.Sort(bytesSlice(bList))
	sort.Sort(bytesSlice(bList2))
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
	sort.Sort(bytesSlice(bList))
	sort.Sort(bytesSlice(bList2))
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
	assert.Equal(t, errors.IsKvNotFound(err), true)
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
