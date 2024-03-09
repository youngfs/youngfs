//go:build rocksdb
// +build rocksdb

package rocksdb

import (
	"dmeta/pkg/kv"
	"github.com/linxGnu/grocksdb"
)

type Iterator struct {
	it     *grocksdb.Iterator
	prefix []byte
}

func (s *Store) NewIterator(opts ...kv.IteratorOption) (kv.Iterator, error) {
	cfg := &kv.IteratorConfig{}
	kv.ApplyConfig(cfg, opts...)

	it := s.db.NewIterator(s.ro)
	return &Iterator{
		it:     it,
		prefix: cfg.Prefix,
	}, nil
}

func (it *Iterator) Seek(key []byte) {
	it.it.Seek(key)
}

func (it *Iterator) Valid() bool {
	if len(it.prefix) == 0 {
		return it.it.Valid()
	}
	return it.it.ValidForPrefix(it.prefix)
}

func (it *Iterator) Next() {
	it.it.Next()
}

func (it *Iterator) Key() []byte {
	key := it.it.Key()
	if key == nil {
		return nil
	}
	data := make([]byte, key.Size())
	copy(data, key.Data())
	return data
}

func (it *Iterator) Value() []byte {
	val := it.it.Value()
	if val == nil {
		return nil
	}
	data := make([]byte, val.Size())
	copy(data, val.Data())
	return data
}

func (it *Iterator) Close() {
	it.it.Close()
}
