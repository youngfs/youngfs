package ec_store

import (
	"context"
	"github.com/go-playground/assert/v2"
	"icesos/command/vars"
	"icesos/kv/redis"
	"icesos/set"
	"icesos/storage_engine/seaweedfs"
	"icesos/util"
	"testing"
)

func TestSetRules_InsertDeleteGetSetRules(t *testing.T) {
	kvStore := redis.NewKvStore(vars.RedisHostPost, vars.RedisPassword, vars.RedisDatabase)
	se := seaweedfs.NewStorageEngine(vars.MasterServer, kvStore)
	client := NewEC(kvStore, se)

	ctx := context.Background()

	hosts, err := se.GetHosts(ctx)
	assert.Equal(t, err, nil)

	setName := set.Set(util.RandString(16))

	setRules := &set.SetRules{
		Set:             setName,
		Hosts:           hosts,
		DataShards:      uint64(len(hosts) - 1),
		ParityShards:    1,
		MAXShardSize:    16 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}

	err = client.InsertSetRules(ctx, setRules)
	assert.Equal(t, err, nil)

	setRules2, err := client.GetSetRules(ctx, setName, true)
	assert.Equal(t, err, nil)
	assert.Equal(t, setRules2, setRules)

	setRules = &set.SetRules{
		Set:             setName,
		Hosts:           hosts,
		DataShards:      uint64(len(hosts) - 1),
		ParityShards:    1,
		MAXShardSize:    16 * 1024 * 1024,
		ECMode:          false,
		ReplicationMode: false,
	}

	err = client.InsertSetRules(ctx, setRules)
	assert.Equal(t, err, nil)

	setRules2, err = client.GetSetRules(ctx, setName, true)
	assert.Equal(t, err, nil)
	assert.Equal(t, setRules2, setRules)

	err = client.DeleteSetRules(ctx, setName, true)
	assert.Equal(t, err, nil)

	err = client.DeleteSetRules(ctx, setName, true)
	assert.Equal(t, err, nil)

	setRules2, err = client.GetSetRules(ctx, setName, true)
	assert.Equal(t, err, nil)
	assert.Equal(t, setRules2, nil)
}
