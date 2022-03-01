package redis_store

import (
	"context"
	"github.com/go-redis/redis/v8"
	"icesos/errors"
	"icesos/kv"
)

func (store *redisStore) KvPut(ctx context.Context, key string, val []byte) error {
	_, err := store.client.Set(ctx, key, val, 0).Result()
	if err != nil {
		return errors.ErrorCodeResponse[errors.ErrKvSever]
	}
	return nil
}

func (store *redisStore) KvGet(ctx context.Context, key string) ([]byte, error) {
	val, err := store.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, kv.KvNotFound
		} else {
			return nil, errors.ErrorCodeResponse[errors.ErrKvSever]
		}
	}
	return []byte(val), nil
}

func (store *redisStore) KvDelete(ctx context.Context, key string) (bool, error) {
	ret, err := store.client.Del(ctx, key).Result()
	if err != nil {
		err = errors.ErrorCodeResponse[errors.ErrKvSever]
	}
	return ret != 0, nil
}
