//go:build rocksdb
// +build rocksdb

package rocksdb

import (
	"dmeta/pkg/kv"
	"github.com/linxGnu/grocksdb"
	"time"
)

type TTLStore struct {
	db *grocksdb.DB
	// option
	ro *grocksdb.ReadOptions
	wo *grocksdb.WriteOptions
}

func (s *TTLStore) optionInit() {
	s.ro = grocksdb.NewDefaultReadOptions()
	s.wo = grocksdb.NewDefaultWriteOptions()
}

// NewTTLStore should be used to open the db when key-values inserted are meant to be removed from the db in a non-strict 'ttl' amount of time therefore,
// this guarantees that key-values inserted will remain in the db for at least ttl amount of time and the db will make efforts to remove the key-values as soon as possible after ttl seconds of their insertion.
func NewTTLStore(path string, ttl time.Duration, opts ...Option) (*TTLStore, error) {
	cfg := &config{}
	for _, opt := range opts {
		opt.apply(cfg)
	}

	dbOpts := grocksdb.NewDefaultOptions()
	dbOpts.SetCreateIfMissing(true)
	dbOpts.SetLevelCompactionDynamicLevelBytes(true)
	bbtOpts := grocksdb.NewDefaultBlockBasedTableOptions()
	bbtOpts.SetFilterPolicy(grocksdb.NewBloomFilterFull(8))
	dbOpts.SetBlockBasedTableFactory(bbtOpts)
	if cfg.compactionFilter != nil {
		dbOpts.SetCompactionFilter(cfg.compactionFilter)
	}
	db, err := grocksdb.OpenDbWithTTL(dbOpts, path, int(ttl.Seconds()))
	if err != nil {
		return nil, err
	}

	ret := &TTLStore{}
	ret.optionInit()
	ret.db = db
	return ret, nil
}

func (s *TTLStore) NewIterator(opts ...kv.IteratorOption) (kv.Iterator, error) {
	cfg := &kv.IteratorConfig{}
	kv.ApplyConfig(cfg, opts...)

	it := s.db.NewIterator(s.ro)
	return &Iterator{
		it:     it,
		prefix: cfg.Prefix,
	}, nil
}

func (s *TTLStore) Close() error {
	s.ro.Destroy()
	s.wo.Destroy()
	s.db.Close()
	return nil
}
