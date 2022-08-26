package redis

import (
	"context"
	"icesfs/command/vars"
	"icesfs/errors"
	"icesfs/kv"
	"icesfs/log"
)

func (store *KvStore) SAdd(ctx context.Context, key string, member []byte) error {
	_, err := store.client.SAdd(ctx, key, member).Result()
	if err != nil {
		log.Errorw("redis sadd error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "key", key)
		err = errors.GetAPIErr(errors.ErrKvSever)
	}
	return err
}

func (store *KvStore) SMembers(ctx context.Context, key string) ([][]byte, error) {
	val, err := store.client.SMembers(ctx, key).Result()
	if err != nil {
		log.Errorw("redis smembers error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "key", key)
		return nil, errors.GetAPIErr(errors.ErrKvSever)
	}
	if len(val) == 0 {
		return [][]byte{}, kv.NotFound
	}

	ret := make([][]byte, len(val))
	for i, str := range val {
		ret[i] = []byte(str)
	}

	return ret, nil
}

func (store *KvStore) SCard(ctx context.Context, key string) (int64, error) {
	ret, err := store.client.SCard(ctx, key).Result()
	if err != nil {
		log.Errorw("redis scard error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "key", key)
		err = errors.GetAPIErr(errors.ErrKvSever)
	}
	return ret, err
}

func (store *KvStore) SRem(ctx context.Context, key string, member []byte) (bool, error) {
	ret, err := store.client.SRem(ctx, key, member).Result()
	if err != nil {
		log.Errorw("redis srem error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "key", key)
		err = errors.GetAPIErr(errors.ErrKvSever)
	}
	return ret != 0, err
}

func (store *KvStore) SIsMember(ctx context.Context, key string, member []byte) (bool, error) {
	ret, err := store.client.SIsMember(ctx, key, member).Result()
	if err != nil {
		log.Errorw("redis sismember error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "key", key)
		return false, errors.GetAPIErr(errors.ErrKvSever)
	}
	return ret, err
}

// delete all members of the set
func (store *KvStore) SDelete(ctx context.Context, key string) (bool, error) {
	cnt, err := store.SCard(ctx, key)
	if err != nil {
		log.Errorw("redis scard error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "key", key)
		return false, errors.GetAPIErr(errors.ErrKvSever)
	}
	if cnt == 0 {
		return false, nil
	}

	_, err = store.client.SPopN(ctx, key, cnt).Result()
	if err != nil {
		log.Errorw("redis spopn error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "key", key)
		return false, errors.GetAPIErr(errors.ErrKvSever)
	}
	return true, nil
}
