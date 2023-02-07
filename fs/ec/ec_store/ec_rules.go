package ec_store

import (
	"context"
	"youngfs/errors"
	"youngfs/fs/rules"
	fs_set "youngfs/fs/set"
)

func rulesKey(set fs_set.Set) string {
	return string(set) + rulesKv
}

func rulesLockKey(set fs_set.Set) string {
	return string(set) + rulesLock
}

func setTurnKey(set fs_set.Set) string {
	return string(set) + setTurnsKv
}

func (ec *ECStore) InsertRules(ctx context.Context, rules *rules.Rules) error {
	mutex := ec.kvStore.NewMutex(rulesLockKey(rules.Set))
	if err := mutex.Lock(); err != nil {
		return errors.ErrRedisSync.WithStack()
	}
	defer func() {
		_, _ = mutex.Unlock()
	}()

	if !rules.Set.IsLegal() {
		return errors.ErrIllegalSetName
	}

	if !rules.IsLegal() {
		return errors.ErrIllegalSetRules
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
	for _, u := range rules.Hosts {
		if !hostSet[u] {
			return errors.ErrIllegalSetRules
		}
	}

	err = ec.DeleteRules(ctx, rules.Set, false)
	if err != nil {
		return err
	}

	proto, err := rules.EncodeProto(ctx)
	if err != nil {
		return err
	}

	err = ec.kvStore.KvPut(ctx, rulesKey(rules.Set), proto)
	if err != nil {
		return err
	}

	err = ec.kvStore.SetNum(ctx, setTurnKey(rules.Set), 0)
	if err != nil {
		return err
	}

	if rules.ECMode {
		err = ec.initPlan(ctx, rules.Set)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ec *ECStore) DeleteRules(ctx context.Context, set fs_set.Set, lock bool) error {
	if lock {
		mutex := ec.kvStore.NewMutex(rulesLockKey(set))
		if err := mutex.Lock(); err != nil {
			return errors.ErrRedisSync.WithStack()
		}
		defer func() {
			_, _ = mutex.Unlock()
		}()
	}

	r, err := ec.GetRules(ctx, set)
	if err != nil {
		if errors.IsKvNotFound(err) {
			return nil
		}
		return err
	}

	// clear set rules map
	ec.rulesMap.Delete(set)

	_, err = ec.kvStore.KvDelete(ctx, rulesKey(set))
	if err != nil && errors.IsKvNotFound(err) {
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

	if r != nil {
		for turns := 0; turns < int(r.DataShards+r.ParityShards); turns++ {
			_, err := ec.kvStore.SDelete(ctx, setPlanShardKey(set, turns))
			if err != nil {
				return errors.WithMessage(err, "delete set rules error: clear plan shard")
			}
		}
	}

	return nil
}

// if don't have set rules, return nil,nil
func (ec *ECStore) GetRules(ctx context.Context, set fs_set.Set) (*rules.Rules, error) {
	if val, ok := ec.rulesMap.Load(set); !ok {
		if r, ok := val.(*rules.Rules); !ok {
			return r, nil
		}
	}

	proto, err := ec.kvStore.KvGet(ctx, rulesKey(set))
	if err != nil {
		return nil, err
	}

	r, err := rules.DecodeRulesProto(ctx, proto)
	if err != nil {
		return nil, err
	}

	ec.rulesMap.Store(set, r)
	return r, nil
}
