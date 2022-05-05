package ec

import (
	"context"
	"icesos/command/vars"
	"icesos/entry"
	"icesos/kv"
	"icesos/log"
	"icesos/set"
	"icesos/storage_engine"
	"icesos/util"
	"strconv"
)

type EC struct {
	kvStore       kv.KvStoreWithRedisMutex
	setRulesMap   map[set.Set]*set.SetRules
	ecQueue       *util.UnboundedQueue[string]
	storageEngine storage_engine.StorageEngine
}

func NewEC(kvStore kv.KvStoreWithRedisMutex, storageEngine storage_engine.StorageEngine) *EC {
	return &EC{
		kvStore:       kvStore,
		setRulesMap:   make(map[set.Set]*set.SetRules),
		ecQueue:       util.NewUnboundedQueue[string](),
		storageEngine: storageEngine,
	}
}

func (ec *EC) genECid(ctx context.Context) (string, error) {
	num, err := ec.kvStore.Incr(ctx, genECidKv)
	if err != nil {
		log.Errorw("gen ECid failed, kv store can't incr", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error())
		return "", err
	}
	return strconv.FormatInt(num, 10), nil
}

// return host, fid, err
func (ec *EC) InsertObject(ent *entry.Entry) (string, string, error) {
	return "", "", nil
}
