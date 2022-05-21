package ec_store

import (
	"context"
	"github.com/go-playground/assert/v2"
	"icesos/command/vars"
	"icesos/entry"
	"icesos/full_path"
	"icesos/kv"
	"icesos/kv/redis"
	"icesos/set"
	"icesos/storage_engine/seaweedfs"
	"icesos/util"
	"testing"
)

func TestEC_RecoverEC(t *testing.T) {
	kvStore := redis.NewKvStore(vars.RedisSocket, vars.RedisPassword, vars.RedisDatabase)
	se := seaweedfs.NewStorageEngine(vars.SeaweedFSMaster, kvStore)
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

	plan, err := client.getPlan(ctx, setName)
	assert.Equal(t, err, nil)

	ent := &entry.Entry{
		Set:      setName,
		FullPath: full_path.FullPath(util.RandString(16)),
		FileSize: 1 * 1024 * 1024,
	}

	_, ecid, _, err := client.InsertObject(ctx, ent)
	assert.Equal(t, err, nil)

	ent.ECid = ecid

	err = client.RecoverEC(ctx, ent)
	assert.Equal(t, err, nil)

	plan2, err := client.getPlan(ctx, setName)
	assert.Equal(t, err, nil)
	assert.Equal(t, plan2, plan)

	suite, err := client.GetSuite(ctx, ecid)
	assert.Equal(t, err, kv.NotFound)
	assert.Equal(t, suite, nil)

	for i := range hosts {
		frags, err := client.getFrags(ctx, setName, i)
		assert.Equal(t, err, nil)
		assert.Equal(t, frags, nil)
	}

	err = client.DeleteSetRules(ctx, setName, true)
	assert.Equal(t, err, nil)
}
