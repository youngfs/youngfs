package set_iam

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestUserIAM_EnDecodeProto(t *testing.T) {
	val := &UserIAM{
		User:      "test",
		SecretKey: "password",
	}

	b, err := val.EncodeProto()
	assert.Equal(t, err, nil)

	val2, err := DecodeUserIAMProto(b)
	assert.Equal(t, err, nil)
	assert.Equal(t, val2, val)
}
