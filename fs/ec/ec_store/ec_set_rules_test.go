package ec_store

import (
	"context"
	"github.com/go-playground/assert/v2"
	"testing"
	fs_set "youngfs/fs/set"
	"youngfs/fs/storage_engine/seaweedfs"
	"youngfs/kv/redis"
	"youngfs/util"
	"youngfs/vars"
)

func TestSetRules_InsertDeleteGetSetRules(t *testing.T) {
	kvStore := redis.NewKvStore(vars.RedisSocket, vars.RedisPassword, vars.RedisDatabase)
	se := seaweedfs.NewStorageEngine(vars.SeaweedFSMaster, kvStore)
	client := NewEC(kvStore, se)

	ctx := context.Background()

	hosts, err := se.GetHosts(ctx)
	assert.Equal(t, err, nil)

	set := fs_set.Set(util.RandString(16))

	setRules := &fs_set.SetRules{
		Set:             set,
		Hosts:           hosts,
		DataShards:      uint64(len(hosts) - 1),
		ParityShards:    1,
		MAXShardSize:    16 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}

	err = client.InsertSetRules(ctx, setRules)
	assert.Equal(t, err, nil)

	setRules2, err := client.GetSetRules(ctx, set, true)
	assert.Equal(t, err, nil)
	assert.Equal(t, setRules2, setRules)

	setRules = &fs_set.SetRules{
		Set:             set,
		Hosts:           hosts,
		DataShards:      uint64(len(hosts) - 1),
		ParityShards:    1,
		MAXShardSize:    16 * 1024 * 1024,
		ECMode:          false,
		ReplicationMode: false,
	}

	err = client.InsertSetRules(ctx, setRules)
	assert.Equal(t, err, nil)

	setRules2, err = client.GetSetRules(ctx, set, true)
	assert.Equal(t, err, nil)
	assert.Equal(t, setRules2, setRules)

	err = client.DeleteSetRules(ctx, set, true)
	assert.Equal(t, err, nil)

	err = client.DeleteSetRules(ctx, set, true)
	assert.Equal(t, err, nil)

	setRules2, err = client.GetSetRules(ctx, set, true)
	assert.Equal(t, err, nil)
	assert.Equal(t, setRules2, nil)
}
