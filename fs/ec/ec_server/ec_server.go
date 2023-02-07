package ec_server

import (
	"context"
	"sync/atomic"
	"time"
	"youngfs/errors"
	"youngfs/fs/ec/ec_calc"
	"youngfs/fs/ec/ec_store"
	"youngfs/fs/entry"
	"youngfs/fs/rules"
	fs_set "youngfs/fs/set"
	"youngfs/util"
)

type ECServer struct {
	ECStore *ec_store.ECStore
	ECCalc  *ec_calc.ECCalc
	ecQueue *util.UnboundedQueue[string]
}

func NewECServer(ecStore *ec_store.ECStore, ecCalc *ec_calc.ECCalc) *ECServer {
	ecServer := &ECServer{
		ECStore: ecStore,
		ECCalc:  ecCalc,
		ecQueue: util.NewUnboundedQueue[string](),
	}

	go ecServer.loopProcessingEC()

	return ecServer
}

func (ec *ECServer) loopProcessingEC() {
	var ecCnt int64
	lce := util.NewLimitedConcurrentExecutor(4)
	for {
		ecCnt = 0
		ec.ecQueue.Consume(func(ecids []string) {
			for _, id := range ecids {
				ecid := id
				lce.Execute(func() {
					ctx := context.Background()
					_ = ec.ECCalc.ExecEC(ctx, ecid)
					atomic.AddInt64(&ecCnt, 1)
				})
			}
			lce.Wait()
		})
		if ecCnt == 0 {
			time.Sleep(time.Second)
		}
	}
}

// return host, ecid, error
func (ec *ECServer) InsertObject(ctx context.Context, ent *entry.Entry) (string, string, error) {
	host, ecid, _, err := ec.ECStore.InsertObject(ctx, ent)
	if err != nil {
		return "", "", err
	}

	return host, ecid, nil
}

func (ec *ECServer) RecoverObject(ctx context.Context, ent *entry.Entry) ([]ec_store.Frag, error) {
	return ec.ECCalc.RecoverObject(ctx, ent)
}

func (ec *ECServer) InsertRules(ctx context.Context, setRules *rules.Rules) error {
	return ec.ECStore.InsertRules(ctx, setRules)
}

func (ec *ECServer) DeleteRules(ctx context.Context, set fs_set.Set) error {
	return ec.ECStore.DeleteRules(ctx, set, true)
}

func (ec *ECServer) GetRules(ctx context.Context, set fs_set.Set) (*rules.Rules, error) {
	setRules, err := ec.ECStore.GetRules(ctx, set)
	if err != nil {
		if errors.IsKvNotFound(err) {
			return nil, errors.ErrRulesNotExist
		}
		return nil, err
	}
	return setRules, err
}

func (ec *ECServer) ExecEC(ctx context.Context, ecid string) error {
	if ecid != "" {
		ec.ecQueue.EnQueue(ecid)
	}
	return nil
}

func (ec *ECServer) ConfirmEC(ctx context.Context, ent *entry.Entry) error {
	err := ec.ECStore.ConfirmEC(ctx, ent)
	if err == nil && ent.ECid != "" {
		ec.ecQueue.EnQueue(ent.ECid)
	}
	return err
}

func (ec *ECServer) RecoverEC(ctx context.Context, ent *entry.Entry) error {
	return ec.ECStore.RecoverEC(ctx, ent)
}

func (ec *ECServer) DeleteObject(ctx context.Context, ecid string) error {
	return ec.ECCalc.DeleteObject(ctx, ecid)
}
