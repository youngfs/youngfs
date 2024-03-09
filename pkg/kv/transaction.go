package kv

import (
	"context"
	"github.com/hashicorp/go-multierror"
)

type Transaction interface {
	Put(ctx context.Context, key []byte, val []byte) error
	Get(ctx context.Context, key []byte) ([]byte, error)
	Delete(ctx context.Context, key []byte) error
	Commit(ctx context.Context) error
	Rollback() error
}

func DoTransaction(store TransactionStore, ctx context.Context, f func(txn Transaction) error) error {
	txn, err := store.NewTransaction()
	if err != nil {
		return err
	}
	err = f(txn)
	if err != nil {
		rerr := txn.Rollback()
		if rerr != nil {
			return multierror.Append(err, rerr)
		}
		return err
	}
	err = txn.Commit(ctx)
	if err != nil {
		rerr := txn.Rollback()
		if rerr != nil {
			return multierror.Append(err, rerr)
		}
		return err
	}
	return nil
}
