package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"youngfs/errors"
)

func (store *KvStore) Incr(ctx context.Context, key string) (int64, error) {
	ret, err := store.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, errors.ErrKvSever.Wrap("redis incr error")
	}
	return ret, nil
}

func (store *KvStore) Decr(ctx context.Context, key string) (int64, error) {
	ret, err := store.client.Decr(ctx, key).Result()
	if err != nil {
		return 0, errors.ErrKvSever.Wrap("redis decr error")
	}
	return ret, nil
}

func (store *KvStore) GetNum(ctx context.Context, key string) (int64, error) {
	val, err := store.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, errors.ErrKvNotFound
		} else {
			return 0, errors.ErrKvSever.Wrap("redis get num: kv get error")
		}
	}

	// val is too big, no parse
	if len(val) > 1024 {
		return 0, errors.ErrKvSever.Wrap("redis get num: parse int")
	}

	ret, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, errors.ErrKvSever.Wrap("redis get num: parse int")
	}

	return ret, nil
}

func (store *KvStore) SetNum(ctx context.Context, key string, num int64) error {
	val := strconv.FormatInt(num, 10)
	_, err := store.client.Set(ctx, key, val, 0).Result()
	if err != nil {
		return errors.ErrKvSever.Wrap("redis bucket num: kv put error")
	}
	return nil
}

func (store *KvStore) ClrNum(ctx context.Context, key string) (bool, error) {
	_, err := store.GetNum(ctx, key)
	if err != nil {
		if err == errors.ErrKvNotFound {
			return false, nil
		}
		return false, err
	}

	ret, err := store.client.Del(ctx, key).Result()
	if err != nil {
		err = errors.ErrKvSever.Wrap("redis clr num: kv delete error")
	}
	return ret != 0, err
}
