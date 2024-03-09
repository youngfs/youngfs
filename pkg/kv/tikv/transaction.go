package tikv

import (
	"context"
	tikverr "github.com/tikv/client-go/v2/error"
	"github.com/tikv/client-go/v2/txnkv/transaction"
	"github.com/youngfs/youngfs/pkg/kv"
)

type Transaction struct {
	txn *transaction.KVTxn
}

func (s *Store) NewTransaction() (kv.Transaction, error) {
	txn, err := s.getTxn()
	if err != nil {
		return nil, err
	}
	return &Transaction{
		txn: txn,
	}, nil
}

func (txn *Transaction) Put(ctx context.Context, key []byte, val []byte) error {
	return txn.txn.Set(key, val)
}

func (txn *Transaction) Get(ctx context.Context, key []byte) ([]byte, error) {
	val, err := txn.txn.Get(ctx, key)
	if tikverr.IsErrNotFound(err) {
		return nil, kv.ErrKeyNotFound
	}
	return val, err
}

func (txn *Transaction) Delete(ctx context.Context, key []byte) error {
	return txn.txn.Delete(key)
}

func (txn *Transaction) Commit(ctx context.Context) error {
	return txn.txn.Commit(ctx)
}

func (txn *Transaction) Rollback() error {
	return txn.txn.Rollback()
}
