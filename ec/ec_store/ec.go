package ec_store

import (
	"context"
	"icesos/command/vars"
	"icesos/entry"
	"icesos/errors"
	"icesos/kv"
	"icesos/log"
	"icesos/set"
	"icesos/storage_engine"
	"math/rand"
	"strconv"
)

type ECStore struct {
	kvStore       kv.KvStoreWithRedisMutex
	setRulesMap   map[set.Set]*set.SetRules
	storageEngine storage_engine.StorageEngine
}

func NewEC(kvStore kv.KvStoreWithRedisMutex, storageEngine storage_engine.StorageEngine) *ECStore {
	return &ECStore{
		kvStore:       kvStore,
		setRulesMap:   make(map[set.Set]*set.SetRules),
		storageEngine: storageEngine,
	}
}

func ecidKey(ecid string) string {
	return ecid + ecidKv
}

func (ec *ECStore) genECid(ctx context.Context) (string, error) {
	num, err := ec.kvStore.Incr(ctx, genECidKv)
	if err != nil {
		log.Errorw("gen ECid failed, kv store can't incr", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error())
		return "", err
	}
	return strconv.FormatInt(num, 10), nil
}

// return host, ecid, suiteid, err
func (ec *ECStore) InsertObject(ctx context.Context, ent *entry.Entry) (string, string, string, error) {
	mutex := ec.kvStore.NewMutex(setRulesLockKey(ent.Set))
	if err := mutex.Lock(); err != nil {
		log.Errorw("get set rules lock error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "set", ent.Set, "entry", ent)
		return "", "", "", errors.GetAPIErr(errors.ErrRedisSync)
	}
	defer func() {
		_, _ = mutex.Unlock()
	}()

	setRules, err := ec.GetSetRules(ctx, ent.Set, false)
	if err != nil {
		return "", "", "", err
	}

	if setRules.ECMode && ent.FileSize > setRules.MAXShardSize {
		return "", "", "", errors.GetAPIErr(errors.ErrIllegalObjectSize)
	}

	// not set rules
	if setRules == nil {
		return "", "", "", nil
	}

	ecid, err := ec.genECid(ctx)
	if err != nil {
		log.Errorw("ec_store insert object error: gen ecid", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "ent", ent)
		return "", "", "", errors.GetAPIErr(errors.ErrServer)
	}

	host := ""
	bakHost := ""
	suiteId := ""
	suite := &Suite{
		ECid:     ecid,
		FullPath: ent.FullPath,
		Set:      ent.Set,
		FileSize: ent.FileSize,
	}

	if setRules.ECMode {
		plan, err := ec.getPlan(ctx, ent.Set)
		if err != nil {
			log.Errorw("ec_store insert object error: get plan", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "ent", ent)
			return "", "", "", errors.GetAPIErr(errors.ErrServer)
		}

		for i, u := range plan.Shards {
			if u.ShardSize >= ent.FileSize {
				host = u.Host
				if i >= len(plan.Shards)-1 {
					log.Errorw("ec_store insert object error: get bakhost error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "ent", ent)
					return "", "", "", errors.GetAPIErr(errors.ErrServer)
				}
				bakHost = plan.Shards[i+1].Host
				plan.Shards[i].ShardSize -= ent.FileSize
				err := ec.kvStore.SAdd(ctx, setPlanShardKey(ent.Set, i), []byte(ecid))
				if err != nil {
					log.Errorw("ec_store insert object error: insert plan shard ecid", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "ent", ent)
					return "", "", "", err
				}
			}
		}

		if host == "" {
			suiteId, err = ec.genECid(ctx)
			if err != nil {
				log.Errorw("ec_store insert object error: gen ecid", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "ent", ent)
				return "", "", "", errors.GetAPIErr(errors.ErrServer)
			}

			shards := make([]Shard, plan.DataShards)
			for i := range shards {
				shards[i].Host = plan.Shards[i].Host
				shards[i].Frags, err = ec.getFrags(ctx, ent.Set, i)
				if err != nil {
					return "", "", "", err
				}

				_, err = ec.kvStore.SDelete(ctx, setPlanShardKey(ent.Set, i))
				if err != nil {
					return "", "", "", err
				}
			}

			err := ec.InsertSuite(ctx, &Suite{
				ECid:   suiteId,
				Shards: shards,
			})
			if err != nil {
				return "", "", "", err
			}

			_, err = ec.kvStore.Incr(ctx, setTurnKey(ent.Set))
			if err != nil {
				return "", "", "", err
			}

			err = ec.initPlan(ctx, ent.Set)
			if err != nil {
				return "", "", "", err
			}

			plan, err = ec.getPlan(ctx, ent.Set) // last variable
			if err != nil {
				log.Errorw("ec_store insert object error: get plan", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "ent", ent)
				return "", "", "", errors.GetAPIErr(errors.ErrServer)
			}

			for i, u := range plan.Shards {
				if u.ShardSize >= ent.FileSize {
					host = u.Host
					if i >= len(plan.Shards)-1 {
						log.Errorw("ec_store insert object error: get bakhost error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "ent", ent)
						return "", "", "", errors.GetAPIErr(errors.ErrServer)
					}
					bakHost = plan.Shards[i+1].Host
					plan.Shards[i].ShardSize -= ent.FileSize
					err := ec.kvStore.SAdd(ctx, setPlanShardKey(ent.Set, i), []byte(ecid))
					if err != nil {
						log.Errorw("ec_store insert object error: insert plan shard ecid", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "ent", ent)
						return "", "", "", err
					}
				}
			}

			if host == "" {
				log.Errorw("ec_store insert object error: not get host", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), "ent", ent)
				return "", "", "", errors.GetAPIErr(errors.ErrServer)
			}
		}

		err = ec.insertPlan(ctx, plan)
		if err != nil {
			return "", "", "", err
		}

		suite.OrigHost = host
		if setRules.ReplicationMode {
			suite.BakHost = bakHost
		}

	} else if setRules.ReplicationMode {
		hostp := rand.Intn(len(setRules.Hosts))
		host = setRules.Hosts[hostp%len(setRules.Hosts)]
		bakHost = setRules.Hosts[(hostp+1)%len(setRules.Hosts)]

		suite.OrigHost = host
		suite.BakHost = bakHost
	} else {
		hostp := rand.Intn(len(setRules.Hosts))
		host = setRules.Hosts[hostp%len(setRules.Hosts)]
		ecid = ""
	}

	err = ec.InsertSuite(ctx, suite)
	if err != nil {
		return "", "", "", err
	}

	return host, ecid, suiteId, nil
}

func (ec *ECStore) ConfirmEC(ctx context.Context, ent *entry.Entry) error {
	if ent.ECid == "" {
		return nil
	}

	suite, err := ec.GetSuite(ctx, ent.ECid)
	if err != nil {
		return err
	}

	suite.OrigFid = ent.Fid
	err = ec.InsertSuite(ctx, suite)
	if err != nil {
		return err
	}

	return nil
}

func (ec *ECStore) RecoverEC(ctx context.Context, ent *entry.Entry) error {
	if ent.ECid == "" {
		return nil
	}

	mutex := ec.kvStore.NewMutex(setRulesLockKey(ent.Set))
	if err := mutex.Lock(); err != nil {
		log.Errorw("get set rules lock error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "set", ent.Set, "entry", ent)
		return errors.GetAPIErr(errors.ErrRedisSync)
	}

	setRules, err := ec.GetSetRules(ctx, ent.Set, false)
	if err != nil {
		return err
	}

	suite, err := ec.GetSuite(ctx, ent.ECid)
	if err != nil {
		return err
	}

	if setRules.ECMode {
		plan, err := ec.getPlan(ctx, ent.Set)
		if err != nil {
			log.Errorw("rocover ec_store: get plan", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "ent", ent)
			return errors.GetAPIErr(errors.ErrServer)
		}

		for i, u := range plan.Shards {
			if u.Host == suite.OrigHost {
				plan.Shards[i].ShardSize += ent.FileSize
				_, err := ec.kvStore.SRem(ctx, setPlanShardKey(ent.Set, i), []byte(ent.ECid))
				if err != nil {
					return err
				}
			}
		}
	}

	err = ec.DeleteSuite(ctx, ent.ECid)
	if err != nil {
		return err
	}

	return nil
}
