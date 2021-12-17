package iam

import (
	"crypto/md5"
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestUserIAM_EnDecodeProto(t *testing.T) {
	val := &userIAM{
		User:      "test",
		SecretKey: md5.Sum([]byte("password")),
	}

	b, err := val.encodeProto()
	assert.Equal(t, err, nil)

	val2, err := decodeUserIAMProto(b)
	assert.Equal(t, err, nil)
	assert.Equal(t, val2, val)
}
