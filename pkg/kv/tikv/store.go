package tikv

import (
	"context"
	"github.com/hashicorp/go-multierror"
	"github.com/pingcap/kvproto/pkg/kvrpcpb"
	"github.com/tikv/client-go/v2/txnkv"
	"github.com/tikv/client-go/v2/txnkv/transaction"
)

type Store struct {
	db                *txnkv.Client
	enable1PC         bool
	enablePessimistic bool
}

func New(addr []string, opts ...Option) (*Store, error) {
	cfg := &config{}
	for _, opt := range opts {
		opt.apply(cfg)
	}

	db, err := txnkv.NewClient(addr, txnkv.WithAPIVersion(kvrpcpb.APIVersion_V2))
	if err != nil {
		return nil, err
	}
	return &Store{
		db:                db,
		enable1PC:         cfg.Enable1PC,
		enablePessimistic: cfg.EnablePessimistic,
	}, nil
}

func (s *Store) getTxn() (*transaction.KVTxn, error) {
	// default optimistic transactions
	txn, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	txn.SetEnable1PC(s.enable1PC)
	txn.SetPessimistic(s.enablePessimistic)
	return txn, nil
}

func (s *Store) update(ctx context.Context, f func(*transaction.KVTxn) error) error {
	txn, err := s.getTxn()
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

func (s *Store) view(f func(*transaction.KVTxn) error) error {
	txn, err := s.getTxn()
	if err != nil {
		return err
	}
	return f(txn)
}

func (s *Store) Close() error {
	return s.db.Close()
}
