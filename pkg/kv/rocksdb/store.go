//go:build rocksdb
// +build rocksdb

package rocksdb

import (
	"github.com/linxGnu/grocksdb"
)

type Store struct {
	db *grocksdb.TransactionDB
	// option
	ro *grocksdb.ReadOptions
	wo *grocksdb.WriteOptions
	to *grocksdb.TransactionOptions
}

func (s *Store) optionInit() {
	s.ro = grocksdb.NewDefaultReadOptions()
	s.wo = grocksdb.NewDefaultWriteOptions()
	s.to = grocksdb.NewDefaultTransactionOptions()
}

// New creates a new Store instance.
// path is a directory path where the database files will be stored.
func New(path string, opts ...Option) (*Store, error) {
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
	txnDBOpts := grocksdb.NewDefaultTransactionDBOptions()
	if cfg.compactionFilter != nil {
		dbOpts.SetCompactionFilter(cfg.compactionFilter)
	}

	db, err := grocksdb.OpenTransactionDb(dbOpts, txnDBOpts, path)
	if err != nil {
		return nil, err
	}

	ret := &Store{}
	ret.optionInit()
	ret.db = db
	return ret, nil
}

func (s *Store) Close() error {
	s.to.Destroy()
	s.wo.Destroy()
	s.to.Destroy()
	s.db.Close()
	return nil
}
