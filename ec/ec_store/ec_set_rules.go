package ec_store

import (
	"context"
	"icesos/command/vars"
	"icesos/errors"
	"icesos/kv"
	"icesos/log"
	"icesos/set"
)

func setRulesKey(set set.Set) string {
	return string(set) + setRulesKv
}

func setRulesLockKey(set set.Set) string {
	return string(set) + setRulesLock
}

func setTurnKey(set set.Set) string {
	return string(set) + setTurnsKv
}

func (ec *ECStore) InsertSetRules(ctx context.Context, setRules *set.SetRules) error {
	mutex := ec.kvStore.NewMutex(setRulesLockKey(setRules.Set))
	if err := mutex.Lock(); err != nil {
		log.Errorw("insert set rules lock error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "set", setRules.Set)
		return errors.GetAPIErr(errors.ErrRedisSync)
	}
	defer func() {
		_, _ = mutex.Unlock()
	}()

	if !setRules.Set.IsLegal() {
		return errors.GetAPIErr(errors.ErrIllegalSetName)
	}

	if !setRules.IsLegal() {
		return errors.GetAPIErr(errors.ErrIllegalSetRules)
	}

	hosts, err := ec.storageEngine.GetHosts(ctx)
	if err != nil {
		return err
	}

	hostSet := make(map[string]bool)
	for _, u := range hosts {
		hostSet[u] = true
	}

	//check hosts is existential
	for _, u := range setRules.Hosts {
		if !hostSet[u] {
			return errors.GetAPIErr(errors.ErrIllegalSetRules)
		}
	}

	// clear set rules map
	ec.setRulesMap[setRules.Set] = nil

	proto, err := setRules.EncodeProto(ctx)
	if err != nil {
		return err
	}

	err = ec.kvStore.KvPut(ctx, setRulesKey(setRules.Set), proto)
	if err != nil {
		return err
	}

	err = ec.kvStore.SetNum(ctx, setTurnKey(setRules.Set), 0)
	if err != nil {
		return err
	}

	_, err = ec.kvStore.KvDelete(ctx, setPlanKey(setRules.Set))
	if err != nil {
		return err
	}

	if setRules.ECMode {
		err = ec.initPlan(ctx, setRules.Set)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ec *ECStore) DeleteSetRules(ctx context.Context, set set.Set) error {
	mutex := ec.kvStore.NewMutex(setRulesLockKey(set))
	if err := mutex.Lock(); err != nil {
		log.Errorw("delete set rules lock error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "set", set)
		return errors.GetAPIErr(errors.ErrRedisSync)
	}
	defer func() {
		_, _ = mutex.Unlock()
	}()

	setRules, err := ec.GetSetRules(ctx, set, false)
	if err != nil {
		return err
	}

	// clear set rules map
	ec.setRulesMap[set] = nil

	_, err = ec.kvStore.KvDelete(ctx, setRulesKey(set))
	if err != nil && err != kv.NotFound {
		return err
	}

	_, err = ec.kvStore.ClrNum(ctx, setTurnKey(set))
	if err != nil {
		return err
	}

	_, err = ec.kvStore.KvDelete(ctx, setPlanKey(set))
	if err != nil {
		return err
	}

	if setRules != nil {
		for turns := 0; turns < int(setRules.DataShards+setRules.ParityShards); turns++ {
			_, err := ec.kvStore.SDelete(ctx, setPlanShardKey(set, turns))
			if err != nil {
				log.Errorw("delete set rules error: clear plan shard", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "set", set, "turns", turns)
				return err
			}
		}
	}

	return nil
}

// if don't have set rules, return nil,nil
func (ec *ECStore) GetSetRules(ctx context.Context, setName set.Set, lock bool) (*set.SetRules, error) {
	if lock {
		mutex := ec.kvStore.NewMutex(setRulesLockKey(setName))
		if err := mutex.Lock(); err != nil {
			log.Errorw("get set rules lock error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "set", setName)
			return nil, errors.GetAPIErr(errors.ErrRedisSync)
		}
		defer func() {
			_, _ = mutex.Unlock()
		}()
	}

	if ec.setRulesMap[setName] != nil {
		return ec.setRulesMap[setName], nil
	}

	proto, err := ec.kvStore.KvGet(ctx, setRulesKey(setName))
	if err == kv.NotFound {
		return nil, nil
	}
	if err != nil {
		return nil, nil
	}

	setRules, err := set.DecodeSetRulesProto(ctx, proto)
	if err != nil {
		return nil, err
	}

	ec.setRulesMap[setName] = setRules
	return setRules, nil
}
