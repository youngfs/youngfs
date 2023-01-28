package set

import (
	"context"
	"github.com/go-playground/assert/v2"
	"testing"
	"youngfs/util"
)

func TestSetRules_EnDecodeProto(t *testing.T) {
	ctx := context.Background()

	val := &SetRules{
		Set:             "test",
		Hosts:           []string{util.RandString(16), util.RandString(16), util.RandString(16), util.RandString(16), util.RandString(16), util.RandString(16)},
		DataShards:      4,
		ParityShards:    2,
		MAXShardSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}

	b, err := val.EncodeProto(ctx)
	assert.Equal(t, err, nil)

	val2, err := DecodeSetRulesProto(ctx, b)
	assert.Equal(t, err, nil)
	assert.Equal(t, val2, val)
}
