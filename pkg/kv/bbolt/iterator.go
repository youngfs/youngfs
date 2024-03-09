package bbolt

import (
	"bytes"
	"github.com/youngfs/youngfs/pkg/kv"
	"go.etcd.io/bbolt"
)

type Iterator struct {
	txn    *bbolt.Tx
	cursor *bbolt.Cursor
	prefix []byte
	key    []byte
	val    []byte
}

func (s *Store) NewIterator(opts ...kv.IteratorOption) (kv.Iterator, error) {
	cfg := &kv.IteratorConfig{}
	kv.ApplyConfig(cfg, opts...)

	txn, err := s.db.Begin(false)
	if err != nil {
		return nil, err
	}
	bucket := txn.Bucket(s.bucket)
	cursor := bucket.Cursor()
	return &Iterator{
		txn:    txn,
		cursor: cursor,
		prefix: cfg.Prefix,
	}, nil
}

func (it *Iterator) Seek(key []byte) {
	it.key, it.val = it.cursor.Seek(key)
}

func (it *Iterator) Valid() bool {
	return it.key != nil && bytes.HasPrefix(it.key, it.prefix)
}

func (it *Iterator) Next() {
	it.key, it.val = it.cursor.Next()
}

func (it *Iterator) Key() []byte {
	if it.key == nil {
		return nil
	}
	key := make([]byte, len(it.key))
	copy(key, it.key)
	return key
}

func (it *Iterator) Value() []byte {
	if it.val == nil {
		return nil
	}
	val := make([]byte, len(it.val))
	copy(val, it.val)
	return val
}

func (it *Iterator) Close() {
	_ = it.txn.Rollback()
}
