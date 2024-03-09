package bbolt

import (
	"context"
	"github.com/youngfs/youngfs/pkg/kv"
	"go.etcd.io/bbolt"
)

func (s *Store) Put(ctx context.Context, key []byte, val []byte) error {
	return s.db.Update(func(txn *bbolt.Tx) error {
		b := make([]byte, len(val))
		copy(b, val)
		return txn.Bucket(s.bucket).Put(key, b)
	})
}

func (s *Store) Get(ctx context.Context, key []byte) ([]byte, error) {
	var val []byte
	return val, s.db.View(func(txn *bbolt.Tx) error {
		b := txn.Bucket(s.bucket).Get(key)
		if b == nil {
			return kv.ErrKeyNotFound
		}
		val = make([]byte, len(b))
		copy(val, b)
		return nil
	})
}

func (s *Store) Delete(ctx context.Context, key []byte) error {
	return s.db.Update(func(txn *bbolt.Tx) error {
		return txn.Bucket(s.bucket).Delete(key)
	})
}
