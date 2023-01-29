package ec_store

import (
	"context"
	"strconv"
	"youngfs/errors"
	"youngfs/fs/set"
)

type PlanShard struct {
	Host      string // host
	ShardSize uint64 // shard size
}

type Plan struct {
	set.Set                //set
	DataShards uint64      // data shards
	Shards     []PlanShard // hosts
}

func setPlanKey(set set.Set) string {
	return string(set) + setPlanKv
}

func setPlanShardKey(set set.Set, turns int) string {
	return string(set) + strconv.Itoa(turns) + setPlanShardsKv
}

func setPlanLockKey(set set.Set) string {
	return string(set) + setPlanLock
}

func (ec *ECStore) insertPlan(ctx context.Context, plan *Plan) error {
	proto, err := plan.EncodeProto(ctx)
	if err != nil {
		return err
	}

	err = ec.kvStore.KvPut(ctx, setPlanKey(plan.Set), proto)
	if err != nil {
		return err
	}

	return nil
}

func (ec *ECStore) getPlan(ctx context.Context, set set.Set) (*Plan, error) {
	proto, err := ec.kvStore.KvGet(ctx, setPlanKey(set))
	if err != nil {
		return nil, err
	}

	plan, err := DecodePlanProto(ctx, proto)
	if err != nil {
		return nil, err
	}

	return plan, nil
}

func (ec *ECStore) initPlan(ctx context.Context, set set.Set) error {
	mutex := ec.kvStore.NewMutex(setPlanLockKey(set))
	if err := mutex.Lock(); err != nil {
		return errors.ErrRedisSync.WithStack()
	}
	defer func() {
		_, _ = mutex.Unlock()
	}()

	setRules, err := ec.GetSetRules(ctx, set, false)
	if err != nil {
		if errors.IsKvNotFound(err) {
			return nil
		}
		return err
	}

	turn, err := ec.kvStore.GetNum(ctx, setTurnKey(set))
	if err != nil {
		return errors.WithMessage(err, "init plan error: get set turns")
	}

	shardsNum := setRules.DataShards + setRules.ParityShards
	st := int(setRules.DataShards) * int(turn)
	shards := make([]PlanShard, shardsNum)
	for i := range shards {
		if i < int(setRules.DataShards) {
			shards[i] = PlanShard{
				Host:      setRules.Hosts[(st+i)%len(setRules.Hosts)],
				ShardSize: setRules.MAXShardSize,
			}
		} else {
			shards[i] = PlanShard{
				Host:      setRules.Hosts[(st+i)%len(setRules.Hosts)],
				ShardSize: 0,
			}
		}

	}

	err = ec.insertPlan(ctx, &Plan{
		Set:        setRules.Set,
		DataShards: setRules.DataShards,
		Shards:     shards,
	})
	if err != nil {
		return err
	}

	for turns := 0; turns < int(shardsNum); turns++ {
		_, err := ec.kvStore.SDelete(ctx, setPlanShardKey(setRules.Set, turns))
		if err != nil {
			return errors.WithMessage(err, "init plan error: clear plan shard")
		}
	}

	return nil
}
