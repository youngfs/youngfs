package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"icesos/errors"
	"icesos/kv"
	"strconv"
)

func (store *KvStore) Incr(ctx context.Context, key string) (int64, error) {
	ret, err := store.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, errors.ErrorCodeResponse[errors.ErrKvSever]
	}
	return ret, nil
}

func (store *KvStore) Decr(ctx context.Context, key string) (int64, error) {
	ret, err := store.client.Decr(ctx, key).Result()
	if err != nil {
		return 0, errors.ErrorCodeResponse[errors.ErrKvSever]
	}
	return ret, nil
}

func (store *KvStore) GetNum(ctx context.Context, key string) (int64, error) {
	val, err := store.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, kv.NotFound
		} else {
			return 0, errors.ErrorCodeResponse[errors.ErrKvSever]
		}
	}

	ret, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, errors.ErrorCodeResponse[errors.ErrKvSever]
	}

	return ret, nil
}

func (store *KvStore) SetNum(ctx context.Context, key string, num int64) error {
	val := strconv.FormatInt(num, 10)
	_, err := store.client.Set(ctx, key, val, 0).Result()
	if err != nil {
		return errors.ErrorCodeResponse[errors.ErrKvSever]
	}
	return nil
}

func (store *KvStore) ClrNum(ctx context.Context, key string) (bool, error) {
	_, err := store.GetNum(ctx, key)
	if err != nil {
		return false, err
	}

	ret, err := store.client.Del(ctx, key).Result()
	if err != nil {
		err = errors.ErrorCodeResponse[errors.ErrKvSever]
	}
	return ret != 0, nil
}
