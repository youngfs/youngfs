package ec_store

import (
	"context"
	"math/rand"
	"sync"
	"youngfs/errors"
	"youngfs/fs/entry"
	"youngfs/fs/id_generator"
	"youngfs/fs/storage_engine"
	"youngfs/kv"
)

type ECStore struct {
	kvStore       kv.KvSetStoreWithRedisMutex
	rulesMap      *sync.Map
	storageEngine storage_engine.StorageEngine
	generator     id_generator.IdGenerator
}

func NewECStore(kvStore kv.KvSetStoreWithRedisMutex, storageEngine storage_engine.StorageEngine, generator id_generator.IdGenerator) *ECStore {
	return &ECStore{
		kvStore:       kvStore,
		rulesMap:      &sync.Map{},
		storageEngine: storageEngine,
		generator:     generator,
	}
}

func ecidKey(ecid string) string {
	return ecid + ecidKv
}

func (ec *ECStore) genECid() (string, error) {
	return ec.generator.Generate()
}

// return host, ecid, suiteid, err
func (ec *ECStore) InsertObject(ctx context.Context, ent *entry.Entry) (string, string, string, error) {
	mutex := ec.kvStore.NewMutex(rulesLockKey(ent.Set))
	if err := mutex.Lock(); err != nil {
		return "", "", "", errors.ErrRedisSync.WithStack()
	}
	defer func() {
		_, _ = mutex.Unlock()
	}()

	setRules, err := ec.GetRules(ctx, ent.Set)
	if err != nil {
		if errors.IsKvNotFound(err) {
			return "", "", "", nil
		}
		return "", "", "", err
	}

	if setRules.ECMode && ent.FileSize > setRules.MaxShardSize {
		return "", "", "", errors.ErrIllegalObjectSize
	}

	// not set rules
	if setRules == nil {
		return "", "", "", nil
	}

	ecid, err := ec.genECid()
	if err != nil {
		return "", "", "", errors.ErrServer.WithMessage("ec_store insert object error gen ecid")
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
			return "", "", "", errors.ErrServer.WithMessage("ec_store insert object error get plan")
		}

		for i, u := range plan.Shards {
			if u.ShardSize >= ent.FileSize {
				host = u.Host
				if i >= len(plan.Shards)-1 {
					return "", "", "", errors.ErrServer.WithMessage("ec_store insert object error get bakhost error")
				}
				bakHost = plan.Shards[i+1].Host
				plan.Shards[i].ShardSize -= ent.FileSize
				err := ec.kvStore.SAdd(ctx, setPlanShardKey(ent.Set, i), []byte(ecid))
				if err != nil {
					return "", "", "", errors.WithMessage(err, "ec_store insert object error insert plan shard ecid")
				}
				break
			}
		}

		if host == "" {
			suiteId, err = ec.genECid()
			if err != nil {
				return "", "", "", errors.ErrServer.WithMessage("ec_store insert object error gen ecid")
			}

			shards := make([]Shard, setRules.DataShards+setRules.ParityShards)
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
				ECid:       suiteId,
				DataShards: setRules.DataShards,
				Shards:     shards,
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
				return "", "", "", errors.ErrServer.WithMessage("ec_store insert object error get plan")
			}

			for i, u := range plan.Shards {
				if u.ShardSize >= ent.FileSize {
					host = u.Host
					if i >= len(plan.Shards)-1 {
						return "", "", "", errors.ErrServer.WithMessage("ec_store insert object error: get bakhost error")
					}
					bakHost = plan.Shards[i+1].Host
					plan.Shards[i].ShardSize -= ent.FileSize
					err := ec.kvStore.SAdd(ctx, setPlanShardKey(ent.Set, i), []byte(ecid))
					if err != nil {
						return "", "", "", errors.WithMessage(err, "ec_store insert object error: insert plan shard ecid")
					}
					break
				}
			}

			if host == "" {
				return "", "", "", errors.ErrServer.WithMessage("ec_store insert object error not get host")
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

	mutex := ec.kvStore.NewMutex(rulesLockKey(ent.Set))
	if err := mutex.Lock(); err != nil {
		return errors.ErrRedisSync.WithStack()
	}
	defer func() {
		_, _ = mutex.Unlock()
	}()

	setRules, err := ec.GetRules(ctx, ent.Set)
	if err != nil {
		if errors.IsKvNotFound(err) {
			return nil
		}
		return err
	}

	suite, err := ec.GetSuite(ctx, ent.ECid)
	if err != nil {
		return err
	}

	if setRules.ECMode {
		plan, err := ec.getPlan(ctx, ent.Set)
		if err != nil {
			return errors.ErrServer.WithMessage("rocover ec_store: get plan")
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

		err = ec.insertPlan(ctx, plan)
		if err != nil {
			return err
		}
	}

	err = ec.DeleteSuite(ctx, ent.ECid)
	if err != nil {
		return err
	}

	return nil
}

func ecidLinkKey(ecid string) string {
	return ecid + ecidLink
}

func (ec *ECStore) SetECidLink(ctx context.Context, ecid string, num int64) error {
	err := ec.kvStore.SetNum(ctx, ecidLinkKey(ecid), num)
	if err != nil {
		return err
	}
	return nil
}

func (ec *ECStore) AddECidLink(ctx context.Context, ecid string) (int64, error) {
	return ec.kvStore.Incr(ctx, ecidLinkKey(ecid))
}

func (ec *ECStore) DelECidLink(ctx context.Context, ecid string) (int64, error) {
	return ec.kvStore.Decr(ctx, ecidLinkKey(ecid))
}

func (ec *ECStore) ClrECidLink(ctx context.Context, ecid string) error {
	_, err := ec.kvStore.ClrNum(ctx, ecidLinkKey(ecid))
	return err
}
