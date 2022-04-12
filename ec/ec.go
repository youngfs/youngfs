package ec

import (
	"context"
	"icesos/kv"
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
		return "", err
	}
	return strconv.FormatInt(num, 10), nil
}
