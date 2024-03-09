package badger

import (
	"context"
	"github.com/dgraph-io/badger/v4"
	"github.com/youngfs/youngfs/pkg/kv"
	"time"
)

func (s *Store) Put(ctx context.Context, key []byte, val []byte) error {
	return s.db.Update(func(txn *badger.Txn) error {
		if s.ttl == nil {
			return txn.Set(key, val)
		} else {
			return txn.SetEntry(badger.NewEntry(key, val).WithTTL(*s.ttl))
		}
	})
}

func (s *Store) Get(ctx context.Context, key []byte) ([]byte, error) {
	var val []byte
	return val, s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return kv.ErrKeyNotFound
			}
			return err
		}
		val, err = item.ValueCopy(nil)
		return err
	})
}

func (s *Store) Delete(ctx context.Context, key []byte) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

func (s *Store) PutWithTTL(ctx context.Context, key []byte, val []byte, ttl time.Duration) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.SetEntry(badger.NewEntry(key, val).WithTTL(ttl))
	})
}
