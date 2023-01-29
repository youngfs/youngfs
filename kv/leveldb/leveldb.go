package leveldb

import (
	"github.com/syndtr/goleveldb/leveldb"
	leveldb_errors "github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"os"
	"youngfs/errors"
	"youngfs/log"
)

type KvStore struct {
	db *leveldb.DB
}

func NewKvStore(dir string) (*KvStore, error) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Errorf("leveldb init error :%s", err.Error())
		return nil, errors.ErrKvSever.WrapErr(err)
	}
	opts := &opt.Options{
		BlockCacheCapacity: 32 * 1024 * 1024,
		WriteBuffer:        16 * 1024 * 1024,
		Filter:             filter.NewBloomFilter(10),
	}

	db, err := leveldb.OpenFile(dir, opts)
	if err != nil {
		if leveldb_errors.IsCorrupted(err) {
			db, err = leveldb.RecoverFile(dir, opts)
		}
		if err != nil {
			log.Errorf("leveldb init error :%s", err.Error())
			return nil, errors.ErrKvSever.WrapErr(err)
		}
	}
	return &KvStore{
		db: db,
	}, nil
}
