package leveldb

import (
	"context"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/youngfs/youngfs/pkg/errors"
	"github.com/youngfs/youngfs/pkg/kv"
)

type Transaction struct {
	txn *leveldb.Transaction
}

func (s *Store) NewTransaction() (kv.Transaction, error) {
	txn, err := s.db.OpenTransaction()
	if err != nil {
		return nil, err
	}
	return &Transaction{
		txn: txn,
	}, nil
}

func (txn *Transaction) Put(ctx context.Context, key []byte, val []byte) error {
	return txn.txn.Put(key, val, nil)
}

func (txn *Transaction) Get(ctx context.Context, key []byte) ([]byte, error) {
	val, err := txn.txn.Get(key, nil)
	if errors.Is(err, leveldb.ErrNotFound) {
		return nil, kv.ErrKeyNotFound
	}
	return val, err
}

func (txn *Transaction) Delete(ctx context.Context, key []byte) error {
	return txn.txn.Delete(key, nil)
}

func (txn *Transaction) Commit(ctx context.Context) error {
	return txn.txn.Commit()
}

func (txn *Transaction) Rollback() error {
	txn.txn.Discard()
	return nil
}
