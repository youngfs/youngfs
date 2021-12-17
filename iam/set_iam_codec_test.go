package set_iam

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestSetIAM_EnDecodeProto(t *testing.T) {
	val := &SetIAM{
		User: "test1",
		Set:  "test2",
	}

	b, err := val.EncodeProto()
	assert.Equal(t, err, nil)

	val2, err := DecodeSetIAMProto(b)
	assert.Equal(t, err, nil)
	assert.Equal(t, val2, val)
}
