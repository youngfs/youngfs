package ec_store

import (
	"context"
	"github.com/go-playground/assert/v2"
	"math/rand"
	"strconv"
	"testing"
	"time"
	"youngfs/fs/full_path"
	"youngfs/fs/set"
	"youngfs/util"
)

func TestSuite_EnDecodeProto(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	ctx := context.Background()

	frags := []*Frag{
		&Frag{
			FullPath: full_path.FullPath(util.RandString(16)),
			Set:      set.Set(util.RandString(16)),
			Fid:      util.RandString(16),
			FileSize: rand.Uint64(),
			OldECid:  util.RandString(16),
		},
		&Frag{
			FullPath: full_path.FullPath(util.RandString(16)),
			Set:      set.Set(util.RandString(16)),
			Fid:      util.RandString(16),
			FileSize: rand.Uint64(),
			OldECid:  util.RandString(16),
		},
		&Frag{
			FullPath: full_path.FullPath(util.RandString(16)),
			Set:      set.Set(util.RandString(16)),
			Fid:      util.RandString(16),
			FileSize: rand.Uint64(),
			OldECid:  util.RandString(16),
		},
	}

	for _, frag := range frags {
		b, err := frag.EncodeProto(ctx)
		assert.Equal(t, err, nil)

		val, err := DecodeFragProto(ctx, b)
		assert.Equal(t, err, nil)
		assert.Equal(t, val, frag)
	}

	shards := []*Shard{
		&Shard{
			Host:  util.RandString(16),
			Frags: []Frag{*frags[0], *frags[1], *frags[2]},
		},
		&Shard{
			Host:  util.RandString(16),
			Frags: []Frag{*frags[0], *frags[1]},
			Md5:   util.RandMd5(),
		},
		&Shard{
			Host:  util.RandString(16),
			Frags: nil,
		},
		&Shard{
			Host: util.RandString(16),
		},
	}

	for _, shard := range shards {
		b, err := shard.EncodeProto(ctx)
		assert.Equal(t, err, nil)

		val, err := DecodeShardProto(ctx, b)
		assert.Equal(t, err, nil)
		assert.Equal(t, val, shard)
	}

	suits := []*Suite{
		&Suite{
			ECid:       strconv.FormatUint(rand.Uint64(), 10),
			FullPath:   full_path.FullPath(util.RandString(16)),
			Set:        set.Set(util.RandString(16)),
			OrigHost:   util.RandString(16),
			OrigFid:    util.RandString(16),
			FileSize:   rand.Uint64(),
			BakHost:    util.RandString(16),
			BakFid:     util.RandString(16),
			Next:       util.RandString(16),
			DataShards: rand.Uint64(),
			Shards:     []Shard{*shards[0], *shards[1]},
		},
		&Suite{
			ECid:       strconv.FormatUint(rand.Uint64(), 10),
			FullPath:   full_path.FullPath(util.RandString(16)),
			Set:        set.Set(util.RandString(16)),
			OrigHost:   util.RandString(16),
			OrigFid:    util.RandString(16),
			FileSize:   rand.Uint64(),
			BakHost:    util.RandString(16),
			BakFid:     util.RandString(16),
			Next:       util.RandString(16),
			DataShards: rand.Uint64(),
			Shards:     []Shard{*shards[0], *shards[1], *shards[2], *shards[3]},
		},
		&Suite{
			ECid:       strconv.FormatUint(rand.Uint64(), 10),
			FullPath:   full_path.FullPath(util.RandString(16)),
			Set:        set.Set(util.RandString(16)),
			OrigHost:   util.RandString(16),
			OrigFid:    util.RandString(16),
			FileSize:   rand.Uint64(),
			BakHost:    util.RandString(16),
			BakFid:     util.RandString(16),
			Next:       "",
			DataShards: rand.Uint64(),
			Shards:     nil,
		},
		&Suite{
			ECid:     strconv.FormatUint(rand.Uint64(), 10),
			FullPath: full_path.FullPath(util.RandString(16)),
			Set:      set.Set(util.RandString(16)),
			OrigHost: util.RandString(16),
			OrigFid:  util.RandString(16),
			FileSize: rand.Uint64(),
			BakHost:  util.RandString(16),
			BakFid:   util.RandString(16),
		},
	}

	for _, suit := range suits {
		b, err := suit.EncodeProto(ctx)
		assert.Equal(t, err, nil)

		val, err := DecodeSuiteProto(ctx, b)
		assert.Equal(t, err, nil)
		assert.Equal(t, val, suit)
	}
}
