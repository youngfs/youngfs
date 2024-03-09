package leveldb

import (
	"context"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/youngfs/youngfs/pkg/errors"
	"github.com/youngfs/youngfs/pkg/kv"
)

func (s *Store) Put(ctx context.Context, key []byte, val []byte) error {
	return s.db.Put(key, val, nil)
}

func (s *Store) Get(ctx context.Context, key []byte) ([]byte, error) {
	val, err := s.db.Get(key, nil)
	if errors.Is(err, leveldb.ErrNotFound) {
		return nil, kv.ErrKeyNotFound
	}
	return val, nil
}

func (s *Store) Delete(ctx context.Context, key []byte) error {
	return s.db.Delete(key, nil)
}
