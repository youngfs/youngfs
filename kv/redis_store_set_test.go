package kv

import (
	"github.com/go-playground/assert/v2"
	"icesos/command/vars"
	"icesos/util"
	"math/rand"
	"sort"
	"testing"
)

func TestRedisStore_Set(t *testing.T) {
	Client.Initialize(vars.RedisHostPost, vars.RedisPassword, vars.RedisDatabase)
	key := "test_redis_set"

	b := util.RandByte(512)

	bList := make([][]byte, 10)
	for i := 0; i < 10; i++ {
		bList[i] = util.RandByte(128 + rand.Int()%128)
	}

	// not add members
	bList2, err := Client.SMembers(key)
	assert.Equal(t, err, nil)
	assert.Equal(t, bList2, nil)

	cnt, err := Client.SCard(key)
	assert.Equal(t, err, nil)
	assert.Equal(t, cnt, int64(0))

	ret, err := Client.SRem(key, b)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, false)

	ret, err = Client.SIsMember(key, b)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, false)

	ret, err = Client.SDelete(key)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, false)

	//add members
	for _, m := range bList {
		err = Client.SAdd(key, m)
		assert.Equal(t, err, nil)
	}

	bList2, err = Client.SMembers(key)
	assert.Equal(t, err, nil)
	sort.Sort(util.BytesSlice(bList))
	sort.Sort(util.BytesSlice(bList2))
	assert.Equal(t, bList2, bList)

	cnt, err = Client.SCard(key)
	assert.Equal(t, err, nil)
	assert.Equal(t, cnt, int64(10))

	ret, err = Client.SRem(key, b)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, false)

	ret, err = Client.SIsMember(key, b)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, false)

	for _, m := range bList {
		ret, err = Client.SIsMember(key, m)
		assert.Equal(t, err, nil)
		assert.Equal(t, ret, true)
	}

	// delete some members
	for i := 7; i < 10; i++ {
		ret, err = Client.SRem(key, bList[i])
		assert.Equal(t, err, nil)
		assert.Equal(t, ret, true)
	}

	bDel := bList[7:]
	bList = bList[:7]

	bList2, err = Client.SMembers(key)
	assert.Equal(t, err, nil)
	sort.Sort(util.BytesSlice(bList))
	sort.Sort(util.BytesSlice(bList2))
	assert.Equal(t, bList2, bList)

	cnt, err = Client.SCard(key)
	assert.Equal(t, err, nil)
	assert.Equal(t, cnt, int64(7))

	ret, err = Client.SRem(key, b)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, false)

	for _, m := range bDel {
		ret, err = Client.SRem(key, m)
		assert.Equal(t, err, nil)
		assert.Equal(t, ret, false)
	}

	ret, err = Client.SIsMember(key, b)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, false)

	for _, m := range bList {
		ret, err = Client.SIsMember(key, m)
		assert.Equal(t, err, nil)
		assert.Equal(t, ret, true)
	}

	for _, m := range bDel {
		ret, err = Client.SIsMember(key, m)
		assert.Equal(t, err, nil)
		assert.Equal(t, ret, false)
	}

	// delete all members
	ret, err = Client.SDelete(key)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, true)

	bList2, err = Client.SMembers(key)
	assert.Equal(t, err, nil)
	assert.Equal(t, bList2, nil)

	cnt, err = Client.SCard(key)
	assert.Equal(t, err, nil)
	assert.Equal(t, cnt, int64(0))

	ret, err = Client.SRem(key, b)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, false)

	for _, m := range bList {
		ret, err = Client.SRem(key, m)
		assert.Equal(t, err, nil)
		assert.Equal(t, ret, false)
	}

	for _, m := range bDel {
		ret, err = Client.SRem(key, m)
		assert.Equal(t, err, nil)
		assert.Equal(t, ret, false)
	}

	ret, err = Client.SIsMember(key, b)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, false)

	for _, m := range bList {
		ret, err = Client.SIsMember(key, m)
		assert.Equal(t, err, nil)
		assert.Equal(t, ret, false)
	}

	for _, m := range bDel {
		ret, err = Client.SIsMember(key, m)
		assert.Equal(t, err, nil)
		assert.Equal(t, ret, false)
	}

	ret, err = Client.SDelete(key)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, false)
}
