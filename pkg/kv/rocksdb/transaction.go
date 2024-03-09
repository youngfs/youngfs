//go:build rocksdb
// +build rocksdb

package rocksdb

import (
	"context"
	"dmeta/pkg/kv"
	"github.com/linxGnu/grocksdb"
)

type Transaction struct {
	txn *grocksdb.Transaction
	ro  *grocksdb.ReadOptions
}

func (s *Store) NewTransaction() (kv.Transaction, error) {
	txn := s.db.TransactionBegin(s.wo, s.to, nil)
	return &Transaction{
		txn: txn,
		ro:  s.ro,
	}, nil
}

func (txn *Transaction) Put(ctx context.Context, key []byte, val []byte) error {
	return txn.txn.Put(key, val)
}

func (txn *Transaction) Get(ctx context.Context, key []byte) ([]byte, error) {
	val, err := txn.txn.Get(txn.ro, key)
	if err != nil {
		return nil, err
	}
	defer val.Free()
	if !val.Exists() {
		return nil, kv.ErrKeyNotFound
	}
	data := make([]byte, val.Size())
	copy(data, val.Data())
	return data, err
}

func (txn *Transaction) Delete(ctx context.Context, key []byte) error {
	return txn.txn.Delete(key)
}

func (txn *Transaction) Commit(ctx context.Context) error {
	return txn.txn.Commit()
}

func (txn *Transaction) Rollback() error {
	return txn.txn.Rollback()
}
