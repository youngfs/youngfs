package iam

import (
	"github.com/go-playground/assert/v2"
	"icesos/command/vars"
	"icesos/kv"
	"testing"
)

func TestUser(t *testing.T) {
	kv.Client.Initialize(vars.RedisHostPost, vars.RedisPassword, vars.RedisDatabase)
	user := User("test_name")
	sk := []string{"test_sk", "test_sk2"}
	set := []string{"test_set1", "test_set2", "test_set3"}

	// add user
	for _, str := range sk {
		ret := user.Identify(str)
		assert.Equal(t, ret, false)
	}

	err := user.CreateUser(sk[0])
	assert.Equal(t, err, nil)

	for i, str := range sk {
		ret := user.Identify(str)
		assert.Equal(t, ret, i == 0)
	}

	err = user.CreateUser(sk[1])
	assert.Equal(t, err, nil)

	for i, str := range sk {
		ret := user.Identify(str)
		assert.Equal(t, ret, i == 1)
	}

	// set
	for _, str := range set {
		ret := user.IsOwnSet(Set(str))
		assert.Equal(t, ret, false)
	}

	for i := 0; i < 2; i++ {
		err = user.AddSet(Set(set[i]))
		assert.Equal(t, err, nil)
	}

	for i, str := range set {
		ret := user.IsOwnSet(Set(str))
		assert.Equal(t, ret, i < 2)
	}

	// delete user
	ret, err := user.DeleteUser()
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, true)

	for _, str := range sk {
		ret = user.Identify(str)
		assert.Equal(t, ret, false)
	}

	for _, str := range set {
		ret = user.IsOwnSet(Set(str))
		assert.Equal(t, ret, false)
	}

	ret, err = user.DeleteUser()
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, false)
}
