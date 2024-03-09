package badger

import (
	"context"
	"github.com/dgraph-io/badger/v4"
	"github.com/youngfs/youngfs/pkg/kv"
)

type Transaction struct {
	txn *badger.Txn
}

func (s *Store) NewTransaction() (kv.Transaction, error) {
	txn := s.db.NewTransaction(true)
	return &Transaction{
		txn: txn,
	}, nil
}

func (txn *Transaction) Put(ctx context.Context, key []byte, val []byte) error {
	return txn.txn.Set(key, val)
}

func (txn *Transaction) Get(ctx context.Context, key []byte) ([]byte, error) {
	item, err := txn.txn.Get(key)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, kv.ErrKeyNotFound
		}
		return nil, err
	}
	return item.ValueCopy(nil)
}

func (txn *Transaction) Delete(ctx context.Context, key []byte) error {
	return txn.txn.Delete(key)
}

func (txn *Transaction) Commit(ctx context.Context) error {
	return txn.txn.Commit()
}

func (txn *Transaction) Rollback() error {
	txn.txn.Discard()
	return nil
}
