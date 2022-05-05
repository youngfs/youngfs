package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"icesos/command/vars"
	"icesos/errors"
	"icesos/kv"
	"icesos/log"
	"strconv"
)

func (store *KvStore) Incr(ctx context.Context, key string) (int64, error) {
	ret, err := store.client.Incr(ctx, key).Result()
	if err != nil {
		log.Errorw("redis incr error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "key", key)
		return 0, errors.GetAPIErr(errors.ErrKvSever)
	}
	return ret, nil
}

func (store *KvStore) Decr(ctx context.Context, key string) (int64, error) {
	ret, err := store.client.Decr(ctx, key).Result()
	if err != nil {
		log.Errorw("redis decr error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "key", key)
		return 0, errors.GetAPIErr(errors.ErrKvSever)
	}
	return ret, nil
}

func (store *KvStore) GetNum(ctx context.Context, key string) (int64, error) {
	val, err := store.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, kv.NotFound
		} else {
			log.Errorw("redis get num: kv get error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "key", key)
			return 0, errors.GetAPIErr(errors.ErrKvSever)
		}
	}

	// val is too big, no parse
	if len(val) > 1024 {
		log.Errorw("", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), "key", key)
		return 0, errors.GetAPIErr(errors.ErrKvSever)
	}

	ret, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		log.Errorw("", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), "key", key) // val is too big, zap-log have bug  error also have val info
		return 0, errors.GetAPIErr(errors.ErrKvSever)
	}

	return ret, nil
}

func (store *KvStore) SetNum(ctx context.Context, key string, num int64) error {
	val := strconv.FormatInt(num, 10)
	_, err := store.client.Set(ctx, key, val, 0).Result()
	if err != nil {
		log.Errorw("redis set num: kv put error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "key", key)
		return errors.GetAPIErr(errors.ErrKvSever)
	}
	return nil
}

func (store *KvStore) ClrNum(ctx context.Context, key string) (bool, error) {
	_, err := store.GetNum(ctx, key)
	if err != nil {
		if err == kv.NotFound {
			return false, nil
		}
		return false, err
	}

	ret, err := store.client.Del(ctx, key).Result()
	if err != nil {
		log.Errorw("redis clr num: kv delete error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "key", key)
		err = errors.GetAPIErr(errors.ErrKvSever)
	}
	return ret != 0, nil
}
