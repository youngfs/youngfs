package leveldb

import (
	"context"
	"github.com/syndtr/goleveldb/leveldb"
	"youngfs/errors"
)

func (store *KvStore) KvPut(ctx context.Context, key string, val []byte) error {
	err := store.db.Put([]byte(key), val, nil)
	if err != nil {
		return errors.ErrKvSever.Wrap("leveldb kv put error")
	}
	return nil
}

func (store *KvStore) KvGet(ctx context.Context, key string) ([]byte, error) {
	val, err := store.db.Get([]byte(key), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, errors.ErrKvNotFound
		} else {
			return nil, errors.ErrKvSever.Wrap("leveldb kv get error")
		}
	}

	return val, nil
}

// KvDelete will return false whether the key exists or not
func (store *KvStore) KvDelete(ctx context.Context, key string) (bool, error) {
	err := store.db.Delete([]byte(key), nil)
	if err != nil {
		return false, errors.ErrKvSever.Wrap("leveldb kv delete error")
	}
	return true, nil
}
