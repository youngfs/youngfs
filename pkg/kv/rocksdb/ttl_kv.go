//go:build rocksdb
// +build rocksdb

package rocksdb

import (
	"context"
	"dmeta/pkg/kv"
)

func (s *TTLStore) Put(ctx context.Context, key []byte, val []byte) error {
	return s.db.Put(s.wo, key, val)
}

func (s *TTLStore) Get(ctx context.Context, key []byte) ([]byte, error) {
	val, err := s.db.Get(s.ro, key)
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

func (s *TTLStore) Delete(ctx context.Context, key []byte) error {
	return s.db.Delete(s.wo, key)
}
