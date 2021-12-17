package iam

import (
	"crypto/md5"
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestUserIAM_EnDecodeProto(t *testing.T) {
	val := &UserIAM{
		User:      "test",
		SecretKey: md5.Sum([]byte("password")),
	}

	b, err := val.EncodeProto()
	assert.Equal(t, err, nil)

	val2, err := DecodeUserIAMProto(b)
	assert.Equal(t, err, nil)
	assert.Equal(t, val2, val)
}
