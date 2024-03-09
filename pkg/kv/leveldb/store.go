package leveldb

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type Store struct {
	db *leveldb.DB
}

func New(path string) (*Store, error) {
	db, err := leveldb.OpenFile(path, &opt.Options{
		BlockCacheCapacity: 32 * 1024 * 1024,         // default value is 8MiB
		WriteBuffer:        16 * 1024 * 1024,         // default value is 4MiB
		Filter:             filter.NewBloomFilter(8), // false positive rate 0.02
	})
	if err != nil {
		return nil, err
	}
	return &Store{
		db: db,
	}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}
