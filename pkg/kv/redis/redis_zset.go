package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/youngfs/youngfs/pkg/errors"
)

func (store *KvStore) ZAdd(ctx context.Context, key, member string) error {
	_, err := store.client.ZAdd(ctx, key, &redis.Z{Score: 0, Member: member}).Result()
	if err != nil {
		err = errors.ErrKvSever.Wrap("redis zadd error")
	}
	return err
}

func (store *KvStore) ZCard(ctx context.Context, key string) (int64, error) {
	ret, err := store.client.ZCard(ctx, key).Result()
	if err != nil {
		err = errors.ErrKvSever.Wrap("redis zcard error")
	}
	return ret, err
}

func (store *KvStore) ZRem(ctx context.Context, key, member string) (bool, error) {
	ret, err := store.client.ZRem(ctx, key, member).Result()
	if err != nil {
		err = errors.ErrKvSever.Wrap("redis zrem error")
	}
	return ret != 0, err
}

func (store *KvStore) ZIsMember(ctx context.Context, key, member string) (bool, error) {
	ret, err := store.client.ZRangeByLex(ctx, key, &redis.ZRangeBy{
		Min:    "[" + member,
		Max:    "[" + member,
		Offset: 0,
		Count:  0,
	}).Result()
	if err != nil {
		err = errors.ErrKvSever.Wrap("redis zismember: zrangebylex error")
	}
	return len(ret) != 0, err
}

// [min , max)
// if min = "" : min = "-"
// if max = "" : max = "+"
func (store *KvStore) ZRangeByLex(ctx context.Context, key, min, max string) ([]string, error) {
	if min == "" {
		min = "-"
	} else {
		min = "[" + min
	}

	if max == "" {
		max = "+"
	} else {
		max = "(" + max
	}

	members, err := store.client.ZRangeByLex(ctx, key,
		&redis.ZRangeBy{
			Min:    min,
			Max:    max,
			Offset: 0,
			Count:  0,
		}).Result()
	if err != nil {
		return nil, errors.ErrKvSever.Wrap("redis zrangebylex error")
	}

	if len(members) == 0 {
		err = errors.ErrKvNotFound
	}

	return members, err
}

// [min , max)
// if min = "" : min = "-"
// if max = "" : max = "+"
func (store *KvStore) ZRemRangeByLex(ctx context.Context, key, min, max string) (bool, error) {
	if min == "" {
		min = "-"
	} else {
		min = "[" + min
	}

	if max == "" {
		max = "+"
	} else {
		max = "(" + max
	}

	cnt, err := store.client.ZRemRangeByLex(ctx, key, min, max).Result()
	if err != nil {
		return false, errors.ErrKvSever.Wrap("redis zremrangebylex error")
	}
	if cnt == 0 {
		return false, nil
	}
	return true, err
}
