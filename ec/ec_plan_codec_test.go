package ec

import (
	"context"
	"github.com/go-playground/assert/v2"
	"icesos/set"
	"icesos/util"
	"math/rand"
	"testing"
	"time"
)

func TestPlan_EnDecodeProto(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	ctx := context.Background()

	shards := []*PlanShard{
		&PlanShard{
			Host:      util.RandString(16),
			ShardSize: rand.Uint64(),
		},
		&PlanShard{
			Host: util.RandString(16),
		},
		&PlanShard{
			ShardSize: rand.Uint64(),
		},
	}

	for _, shard := range shards {
		b, err := shard.EncodeProto(ctx)
		assert.Equal(t, err, nil)

		val, err := DecodePlanShardProto(ctx, b)
		assert.Equal(t, err, nil)
		assert.Equal(t, val, shard)
	}

	plan := &Plan{
		Set: set.Set(util.RandString(16)),
		// shard = nil
	}

	b, err := plan.EncodeProto(ctx)
	assert.Equal(t, err, nil)

	val, err := DecodePlanProto(ctx, b)
	assert.Equal(t, err, nil)
	assert.Equal(t, val, plan)

	plan = &Plan{
		Set:        set.Set(util.RandString(16)),
		DataShards: uint64(len(shards)),
		Shards:     []PlanShard{*shards[0], *shards[1], *shards[2]},
	}

	b, err = plan.EncodeProto(ctx)
	assert.Equal(t, err, nil)

	val, err = DecodePlanProto(ctx, b)
	assert.Equal(t, err, nil)
	assert.Equal(t, val, plan)
}
