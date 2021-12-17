package kv

import (
	"context"
)

func (store *redisStore) KvPut(key string, val []byte) error {
	_, err := store.client.Set(context.Background(), key, val, 0).Result()
	return err
}

func (store *redisStore) KvGet(key string) ([]byte, error) {
	val, err := store.client.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	return []byte(val), err
}

func (store *redisStore) KvDelete(key string) (bool, error) {
	ret, err := store.client.Del(context.Background(), key).Result()
	return ret != 0, err
}
