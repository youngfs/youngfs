package iam

import (
	"github.com/go-playground/assert/v2"
	"icesos/command/vars"
	"icesos/kv/redis"
	"testing"
)

func TestUser(t *testing.T) {
	redis.Client.Initialize(vars.RedisHostPost, vars.RedisPassword, vars.RedisDatabase)
	user := User("test_name")
	sk := []string{"test_sk", "test_sk2"}
	set := []string{"test_set1", "test_set2", "test_set3"}

	// add user
	ret, err := user.IsExist()
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, false)

	for _, str := range sk {
		ret := user.Identify(str)
		assert.Equal(t, ret, false)
	}

	err = user.Create(sk[0])
	assert.Equal(t, err, nil)

	ret, err = user.IsExist()
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, true)

	for i, str := range sk {
		ret := user.Identify(str)
		assert.Equal(t, ret, i == 0)
	}

	err = user.Create(sk[1])
	assert.Equal(t, err, nil)

	ret, err = user.IsExist()
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, true)

	for i, str := range sk {
		ret := user.Identify(str)
		assert.Equal(t, ret, i == 1)
	}

	// set
	for _, str := range set {
		ret := user.ReadSetPermission(Set(str))
		assert.Equal(t, ret, false)
		ret = user.WriteSetPermission(Set(str))
		assert.Equal(t, ret, false)
	}

	for i := 0; i < 2; i++ {
		err = user.AddReadSetPermission(Set(set[i]))
		assert.Equal(t, err, nil)
		err = user.AddReadSetPermission(Set(set[i]))
		assert.Equal(t, err, nil)
	}

	for i := 1; i < 3; i++ {
		err = user.AddWriteSetPermission(Set(set[i]))
		assert.Equal(t, err, nil)
		err = user.AddWriteSetPermission(Set(set[i]))
		assert.Equal(t, err, nil)
	}

	for i, str := range set {
		ret := user.ReadSetPermission(Set(str))
		assert.Equal(t, ret, i < 2)
		ret = user.WriteSetPermission(Set(str))
		assert.Equal(t, ret, i > 0)
	}

	for _, str := range set {
		err := user.DeleteReadSetPermission(Set(str))
		assert.Equal(t, err, nil)
		err = user.DeleteWriteSetPermission(Set(str))
		assert.Equal(t, err, nil)
	}

	for _, str := range set {
		ret := user.ReadSetPermission(Set(str))
		assert.Equal(t, ret, false)
		ret = user.WriteSetPermission(Set(str))
		assert.Equal(t, ret, false)
	}

	ret, err = user.IsExist()
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, true)

	// delete user
	ret, err = user.Delete()
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, true)

	for _, str := range sk {
		ret = user.Identify(str)
		assert.Equal(t, ret, false)
	}

	for _, str := range set {
		ret := user.ReadSetPermission(Set(str))
		assert.Equal(t, ret, false)
		ret = user.WriteSetPermission(Set(str))
		assert.Equal(t, ret, false)
	}

	ret, err = user.Delete()
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, false)
}
