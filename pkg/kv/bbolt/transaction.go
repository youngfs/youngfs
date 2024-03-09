package bbolt

import (
	"context"
	"github.com/youngfs/youngfs/pkg/kv"
	"go.etcd.io/bbolt"
)

type Transaction struct {
	txn    *bbolt.Tx
	bucket *bbolt.Bucket
}

func (s *Store) NewTransaction() (kv.Transaction, error) {
	txn, err := s.db.Begin(true)
	if err != nil {
		return nil, err
	}
	bucket := txn.Bucket(s.bucket)
	return &Transaction{
		txn:    txn,
		bucket: bucket,
	}, nil
}

func (txn *Transaction) Put(ctx context.Context, key []byte, val []byte) error {
	b := make([]byte, len(val))
	copy(b, val)
	return txn.bucket.Put(key, b)
}

func (txn *Transaction) Get(ctx context.Context, key []byte) ([]byte, error) {
	val := txn.bucket.Get(key)
	if val == nil {
		return nil, kv.ErrKeyNotFound
	}
	ret := make([]byte, len(val))
	copy(ret, val)
	return ret, nil
}

func (txn *Transaction) Delete(ctx context.Context, key []byte) error {
	return txn.bucket.Delete(key)
}

func (txn *Transaction) Commit(ctx context.Context) error {
	return txn.txn.Commit()
}

func (txn *Transaction) Rollback() error {
	return txn.txn.Rollback()
}
