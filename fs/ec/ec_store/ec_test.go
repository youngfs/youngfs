package ec_store

import (
	"context"
	"github.com/go-playground/assert/v2"
	"testing"
	"youngfs/errors"
	"youngfs/fs/entry"
	"youngfs/fs/full_path"
	fs_set "youngfs/fs/set"
	"youngfs/fs/storage_engine/seaweedfs"
	"youngfs/kv/redis"
	"youngfs/util"
	"youngfs/vars"
)

func TestEC_RecoverEC(t *testing.T) {
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

	plan, err := client.getPlan(ctx, set)
	assert.Equal(t, err, nil)

	ent := &entry.Entry{
		Set:      set,
		FullPath: full_path.FullPath(util.RandString(16)),
		FileSize: 1 * 1024 * 1024,
	}

	_, ecid, _, err := client.InsertObject(ctx, ent)
	assert.Equal(t, err, nil)

	ent.ECid = ecid

	err = client.RecoverEC(ctx, ent)
	assert.Equal(t, err, nil)

	plan2, err := client.getPlan(ctx, set)
	assert.Equal(t, err, nil)
	assert.Equal(t, plan2, plan)

	suite, err := client.GetSuite(ctx, ecid)
	assert.Equal(t, errors.IsKvNotFound(err), true)
	assert.Equal(t, suite, nil)

	for i := range hosts {
		frags, err := client.getFrags(ctx, set, i)
		assert.Equal(t, err, nil)
		assert.Equal(t, frags, nil)
	}

	err = client.DeleteSetRules(ctx, set, true)
	assert.Equal(t, err, nil)
}
