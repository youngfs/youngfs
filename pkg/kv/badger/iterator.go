package badger

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/youngfs/youngfs/pkg/kv"
)

type Iterator struct {
	txn *badger.Txn
	it  *badger.Iterator
}

func (s *Store) NewIterator(opts ...kv.IteratorOption) (kv.Iterator, error) {
	cfg := &kv.IteratorConfig{}
	kv.ApplyConfig(cfg, opts...)

	txn := s.db.NewTransaction(false)
	itOpts := badger.DefaultIteratorOptions
	itOpts.Prefix = cfg.Prefix
	it := txn.NewIterator(itOpts)
	return &Iterator{
		txn: txn,
		it:  it,
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
	return it.it.Item().KeyCopy(nil)
}

func (it *Iterator) Value() []byte {
	val, _ := it.it.Item().ValueCopy(nil)
	return val
}

func (it *Iterator) Close() {
	it.it.Close()
	it.txn.Discard()
}
