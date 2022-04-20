package ec

import (
	"context"
	"icesos/command/vars"
	"icesos/kv"
	"icesos/log"
	"icesos/set"
	"icesos/util"
	"strconv"
)

type EC struct {
	kvStore     kv.KvStoreWithRedisMutex
	setRulesMap map[set.Set]SetRules
	ecQueue     *util.UnboundedQueue[string]
}

func NewEC(kvStore kv.KvStoreWithRedisMutex) *EC {
	return &EC{
		kvStore:     kvStore,
		setRulesMap: make(map[set.Set]SetRules),
		ecQueue:     util.NewUnboundedQueue[string](),
	}
}

func (ec *EC) genECid(ctx context.Context) (string, error) {
	num, err := ec.kvStore.Incr(ctx, genECidKey)
	if err != nil {
		log.Errorw("gen ECid failed, kv store can't incr", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error())
		return "", err
	}
	return strconv.FormatInt(num, 10), nil
}
