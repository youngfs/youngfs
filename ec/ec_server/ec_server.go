package ec_server

import (
	"context"
	"icesfs/ec/ec_calc"
	"icesfs/ec/ec_store"
	"icesfs/entry"
	"icesfs/errors"
	"icesfs/kv"
	"icesfs/set"
	"icesfs/util"
	"time"
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
	var ecCnt int
	for {
		ecCnt = 0
		ec.ecQueue.Consume(func(ecids []string) {
			for _, ecid := range ecids {
				ctx := context.Background()
				_ = ec.ECCalc.ExecEC(ctx, ecid)
				ecCnt++
			}
		})
		if ecCnt == 0 {
			time.Sleep(1234 * time.Millisecond)
		}
	}
}

// return host, ecid, error
func (ec *ECServer) InsertObject(ctx context.Context, ent *entry.Entry) (string, string, error) {
	host, ecid, suiteid, err := ec.ECStore.InsertObject(ctx, ent)
	if err != nil {
		return "", "", err
	}

	if suiteid != "" {
		ec.ecQueue.EnQueue(suiteid)
	}

	return host, ecid, nil
}

func (ec *ECServer) RecoverObject(ctx context.Context, ent *entry.Entry) ([]ec_store.Frag, error) {
	return ec.ECCalc.RecoverObject(ctx, ent)
}

func (ec *ECServer) InsertSetRules(ctx context.Context, setRules *set.SetRules) error {
	return ec.ECStore.InsertSetRules(ctx, setRules)
}

func (ec *ECServer) DeleteSetRules(ctx context.Context, set set.Set) error {
	return ec.ECStore.DeleteSetRules(ctx, set, true)
}

func (ec *ECServer) GetSetRules(ctx context.Context, setName set.Set) (*set.SetRules, error) {
	setRules, err := ec.ECStore.GetSetRules(ctx, setName, true)
	if err != nil {
		if err == kv.NotFound {
			return nil, errors.GetAPIErr(errors.ErrSetRulesNotExist)
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
	return ec.ECStore.ConfirmEC(ctx, ent)
}

func (ec *ECServer) RecoverEC(ctx context.Context, ent *entry.Entry) error {
	return ec.ECStore.RecoverEC(ctx, ent)
}

func (ec *ECServer) DeleteObject(ctx context.Context, ecid string) error {
	return ec.ECCalc.DeleteObject(ctx, ecid)
}
