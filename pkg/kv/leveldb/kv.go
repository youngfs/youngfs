package leveldb

import (
	"context"
)

func (s *Store) Put(ctx context.Context, key []byte, val []byte) error {
	return s.db.Put(key, val, nil)
}

func (s *Store) Get(ctx context.Context, key []byte) ([]byte, error) {
	return s.db.Get(key, nil)
}

func (s *Store) Delete(ctx context.Context, key []byte) error {
	return s.db.Delete(key, nil)
}
