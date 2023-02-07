package ec_store

import (
	"context"
	"github.com/go-playground/assert/v2"
	"testing"
	"youngfs/fs/id_generator/snow_flake"
	"youngfs/fs/rules"
	fs_set "youngfs/fs/set"
	"youngfs/fs/storage_engine/seaweedfs"
	"youngfs/kv/redis"
	"youngfs/util"
	"youngfs/vars"
)

func TestSetRules_InsertDeleteGetSetRules(t *testing.T) {
	kvStore := redis.NewKvStore(vars.RedisSocket, vars.RedisPassword, vars.RedisDatabase)
	se := seaweedfs.NewStorageEngine(vars.SeaweedFSMaster, kvStore)
	client := NewECStore(kvStore, se, snow_flake.NewSnowFlake(0))

	ctx := context.Background()

	hosts, err := se.GetHosts(ctx)
	assert.Equal(t, err, nil)

	set := fs_set.Set(util.RandString(16))

	setRules := &rules.Rules{
		Set:             set,
		Hosts:           hosts,
		DataShards:      uint64(len(hosts) - 1),
		ParityShards:    1,
		MaxShardSize:    16 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}

	err = client.InsertRules(ctx, setRules)
	assert.Equal(t, err, nil)

	setRules2, err := client.GetRules(ctx, set)
	assert.Equal(t, err, nil)
	assert.Equal(t, setRules2, setRules)

	setRules = &rules.Rules{
		Set:             set,
		Hosts:           hosts,
		DataShards:      uint64(len(hosts) - 1),
		ParityShards:    1,
		MaxShardSize:    16 * 1024 * 1024,
		ECMode:          false,
		ReplicationMode: false,
	}

	err = client.InsertRules(ctx, setRules)
	assert.Equal(t, err, nil)

	setRules2, err = client.GetRules(ctx, set)
	assert.Equal(t, err, nil)
	assert.Equal(t, setRules2, setRules)

	err = client.DeleteRules(ctx, set, true)
	assert.Equal(t, err, nil)

	err = client.DeleteRules(ctx, set, true)
	assert.Equal(t, err, nil)

	setRules2, err = client.GetRules(ctx, set)
	assert.Equal(t, err, nil)
	assert.Equal(t, setRules2, nil)
}
