package ec

import (
	"github.com/go-playground/assert/v2"
	"icesos/util"
	"testing"
)

func TestSetRules_EnDecodeProto(t *testing.T) {
	val := &SetRules{
		Set:             "test",
		Hosts:           []string{util.RandString(16), util.RandString(16), util.RandString(16), util.RandString(16), util.RandString(16), util.RandString(16)},
		DataShards:      4,
		ParityShards:    2,
		MAXBlockSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}

	b, err := val.EncodeProto()
	assert.Equal(t, err, nil)

	val2, err := DecodeSetRulesProto(b)
	assert.Equal(t, err, nil)
	assert.Equal(t, val2, val)
}
