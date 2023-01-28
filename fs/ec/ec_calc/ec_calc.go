package ec_calc

import (
	"context"
	"youngfs/errors"
	ec_store2 "youngfs/fs/ec/ec_store"
	"youngfs/fs/entry"
	"youngfs/fs/storage_engine"
)

type ECCalc struct {
	ECStore       *ec_store2.ECStore
	storageEngine storage_engine.StorageEngine
}

func NewECCalc(ecStore *ec_store2.ECStore, storageEngine storage_engine.StorageEngine) *ECCalc {
	return &ECCalc{
		ECStore:       ecStore,
		storageEngine: storageEngine,
	}
}

func (calc *ECCalc) ExecEC(ctx context.Context, ecid string) error {
	suite, err := calc.ECStore.GetSuite(ctx, ecid)
	if err != nil {
		return err
	}

	for suite.Next != "" {
		suite, err = calc.ECStore.GetSuite(ctx, ecid)
		if err != nil {
			return err
		}
	}

	if suite.BakHost != "" && suite.OrigFid != "" {
		err := calc.backup(ctx, suite)
		if err != nil {
			return err
		}
		return nil
	} else if suite.Shards != nil {
		err := calc.reedSolomon(ctx, suite)
		if err != nil {
			return err
		}
	}

	return nil
}

func (calc *ECCalc) RecoverObject(ctx context.Context, ent *entry.Entry) ([]ec_store2.Frag, error) {
	if ent.ECid == "" {
		return nil, errors.ErrRecoverFailed
	}

	suite := &ec_store2.Suite{}
	suite, err := calc.ECStore.GetSuite(ctx, ent.ECid)
	if err != nil {
		return nil, err
	}
	for suite.Next != "" {
		suite, err = calc.ECStore.GetSuite(ctx, suite.Next)
		if err != nil {
			return nil, err
		}
	}

	if suite.BakFid != "" {
		return calc.backupRecover(ctx, suite, ent)
	} else if suite.Shards != nil {
		return calc.reedSolomonRecover(ctx, suite)
	}

	return nil, errors.ErrRecoverFailed
}

func (calc *ECCalc) DeleteObject(ctx context.Context, ecid string) error {
	if ecid == "" {
		return nil
	}

	suite, err := calc.ECStore.GetSuite(ctx, ecid)
	if err != nil {
		return nil
	}
	for suite.Next != "" {
		suite, err = calc.ECStore.GetSuite(ctx, suite.Next)
		if err != nil {
			return err
		}
	}

	if suite.BakFid != "" {
		err := calc.storageEngine.DeleteObject(ctx, suite.BakFid)
		if err != nil {
			return err
		}
	} else if suite.Shards != nil {
		err := calc.reedSolomonDelete(ctx, suite.ECid)
		if err != nil {
			return err
		}
	}

	err = calc.ECStore.DeleteSuite(ctx, ecid)
	if err != nil {
		return err
	}

	return nil
}
