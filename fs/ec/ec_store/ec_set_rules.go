package ec_store

import (
	"context"
	"youngfs/errors"
	fs_set "youngfs/fs/set"
)

func setRulesKey(set fs_set.Set) string {
	return string(set) + setRulesKv
}

func setRulesLockKey(set fs_set.Set) string {
	return string(set) + setRulesLock
}

func setTurnKey(set fs_set.Set) string {
	return string(set) + setTurnsKv
}

func (ec *ECStore) InsertSetRules(ctx context.Context, setRules *fs_set.SetRules) error {
	mutex := ec.kvStore.NewMutex(setRulesLockKey(setRules.Set))
	if err := mutex.Lock(); err != nil {
		return errors.ErrRedisSync.WithStack()
	}
	defer func() {
		_, _ = mutex.Unlock()
	}()

	if !setRules.Set.IsLegal() {
		return errors.ErrIllegalSetName
	}

	if !setRules.IsLegal() {
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
	for _, u := range setRules.Hosts {
		if !hostSet[u] {
			return errors.ErrIllegalSetRules
		}
	}

	err = ec.DeleteSetRules(ctx, setRules.Set, false)
	if err != nil {
		return err
	}

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

	if setRules.ECMode {
		err = ec.initPlan(ctx, setRules.Set)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ec *ECStore) DeleteSetRules(ctx context.Context, set fs_set.Set, lock bool) error {
	if lock {
		mutex := ec.kvStore.NewMutex(setRulesLockKey(set))
		if err := mutex.Lock(); err != nil {
			return errors.ErrRedisSync.WithStack()
		}
		defer func() {
			_, _ = mutex.Unlock()
		}()
	}

	setRules, err := ec.GetSetRules(ctx, set, false)
	if err != nil {
		if errors.IsKvNotFound(err) {
			return nil
		}
		return err
	}

	// clear set rules map
	ec.setRulesMap[set] = nil

	_, err = ec.kvStore.KvDelete(ctx, setRulesKey(set))
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

	if setRules != nil {
		for turns := 0; turns < int(setRules.DataShards+setRules.ParityShards); turns++ {
			_, err := ec.kvStore.SDelete(ctx, setPlanShardKey(set, turns))
			if err != nil {
				return errors.WithMessage(err, "delete set rules error clear plan shard")
			}
		}
	}

	return nil
}

// if don't have set rules, return nil,nil
func (ec *ECStore) GetSetRules(ctx context.Context, set fs_set.Set, lock bool) (*fs_set.SetRules, error) {
	if lock {
		mutex := ec.kvStore.NewMutex(setRulesLockKey(set))
		if err := mutex.Lock(); err != nil {
			return nil, errors.ErrRedisSync.WithStack()
		}
		defer func() {
			_, _ = mutex.Unlock()
		}()
	}

	if ec.setRulesMap[set] != nil {
		return ec.setRulesMap[set], nil
	}

	proto, err := ec.kvStore.KvGet(ctx, setRulesKey(set))
	if err != nil {
		if errors.IsKvNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	setRules, err := fs_set.DecodeSetRulesProto(ctx, proto)
	if err != nil {
		return nil, err
	}

	ec.setRulesMap[set] = setRules
	return setRules, nil
}
