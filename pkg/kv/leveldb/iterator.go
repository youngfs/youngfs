package leveldb

import (
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/youngfs/youngfs/pkg/kv"
)

type Iterator struct {
	it iterator.Iterator
}

func (s *Store) NewIterator(opts ...kv.IteratorOption) (kv.Iterator, error) {
	cfg := &kv.IteratorConfig{}
	kv.ApplyConfig(cfg, opts...)
	var prefix *util.Range = nil
	if cfg.Prefix != nil {
		prefix = util.BytesPrefix(cfg.Prefix)
	}
	return &Iterator{
		it: s.db.NewIterator(prefix, nil),
	}, nil
}

func (it *Iterator) Seek(key []byte) {
	it.it.Seek(key)
}

func (it *Iterator) Valid() bool {
	return it.it.Valid()
}

func (it *Iterator) Next() {
	it.it.Next()
}

func (it *Iterator) Key() []byte {
	return it.it.Key()
}

func (it *Iterator) Value() []byte {
	return it.it.Value()
}

func (it *Iterator) Close() {
	it.it.Release()
}
