package kv

import (
	"context"
	"github.com/go-redis/redis/v8"
	"icesos/errors"
)

func (store *redisStore) KvPut(key string, val []byte) error {
	_, err := store.client.Set(context.Background(), key, val, 0).Result()
	if err != nil {
		return errors.ErrorCodeResponse[errors.ErrKvSever]
	}
	return nil
}

func (store *redisStore) KvGet(key string) ([]byte, error) {
	val, err := store.client.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, err
		} else {
			return nil, errors.ErrorCodeResponse[errors.ErrKvSever]
		}
	}
	return []byte(val), nil
}

func (store *redisStore) KvDelete(key string) (bool, error) {
	ret, err := store.client.Del(context.Background(), key).Result()
	if err != nil {
		err = errors.ErrorCodeResponse[errors.ErrKvSever]
	}
	return ret != 0, nil
}
