package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"icesos/errors"
	"icesos/kv"
)

func (store *KvStore) KvPut(ctx context.Context, key string, val []byte) error {
	_, err := store.client.Set(ctx, key, val, 0).Result()
	if err != nil {
		return errors.ErrorCodeResponse[errors.ErrKvSever]
	}
	return nil
}

func (store *KvStore) KvGet(ctx context.Context, key string) ([]byte, error) {
	val, err := store.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, kv.NotFound
		} else {
			return nil, errors.ErrorCodeResponse[errors.ErrKvSever]
		}
	}
	return []byte(val), nil
}

func (store *KvStore) KvDelete(ctx context.Context, key string) (bool, error) {
	ret, err := store.client.Del(ctx, key).Result()
	if err != nil {
		err = errors.ErrorCodeResponse[errors.ErrKvSever]
	}
	return ret != 0, nil
}
