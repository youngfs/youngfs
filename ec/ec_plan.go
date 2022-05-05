package ec

import (
	"context"
	"icesos/command/vars"
	"icesos/errors"
	"icesos/log"
	"icesos/set"
	"strconv"
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

func (ec *EC) initPlan(ctx context.Context, setRules *set.SetRules) error {
	mutex := ec.kvStore.NewMutex(setPlanLockKey(setRules.Set))
	if err := mutex.Lock(); err != nil {
		log.Errorw("init plan lock error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "set", setRules.Set)
		return errors.GetAPIErr(errors.ErrRedisSync)
	}
	defer func() {
		_, _ = mutex.Unlock()
	}()

	turn, err := ec.kvStore.GetNum(ctx, setTurnKey(setRules.Set))
	if err != nil {
		log.Errorw("init plan error: get set turns", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "set", setRules.Set)
		return err
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

	plan := &Plan{
		Set:        setRules.Set,
		DataShards: setRules.DataShards,
		Shards:     shards,
	}

	proto, err := plan.EncodeProto(ctx)
	if err != nil {
		log.Errorw("init plan error: encode proto plan", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "set", setRules.Set)
		return err
	}

	err = ec.kvStore.KvPut(ctx, setPlanKey(setRules.Set), proto)
	if err != nil {
		log.Errorw("init plan error: kv put plan", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "set", setRules.Set)
		return err
	}

	for turns := 0; turns < int(shardsNum); turns++ {
		_, err := ec.kvStore.SDelete(ctx, setPlanShardKey(setRules.Set, turns))
		if err != nil {
			log.Errorw("init plan error: clear plan shard", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "set", setRules.Set, "turns", turns)
			return err
		}
	}

	return nil
}
