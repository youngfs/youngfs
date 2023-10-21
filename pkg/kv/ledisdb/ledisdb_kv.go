package ledisdb

import (
	"context"
	"github.com/youngfs/youngfs/pkg/errors"
)

func (store *KvStore) KvPut(ctx context.Context, key string, val []byte) error {
	err := store.db.Set([]byte(key), val)
	if err != nil {
		return errors.ErrKvSever.Wrap("redis kv put error")
	}
	return nil
}

func (store *KvStore) KvGet(ctx context.Context, key string) ([]byte, error) {
	val, err := store.db.Get([]byte(key))
	if err != nil {
		return nil, errors.ErrKvSever.Wrap("redis kv get error")
	}
	if val == nil {
		return nil, errors.ErrKvNotFound
	}
	return val, nil
}

func (store *KvStore) KvDelete(ctx context.Context, key string) (bool, error) {
	ret, err := store.db.Del([]byte(key))
	if err != nil {
		err = errors.ErrKvSever.Wrap("redis kv delete error")
	}
	return ret != 0, nil
}
