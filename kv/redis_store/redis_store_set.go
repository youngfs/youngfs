package redis_store

import (
	"context"
	"icesos/errors"
	"icesos/kv"
)

func (store *RedisStore) SAdd(ctx context.Context, key string, member []byte) error {
	_, err := store.client.SAdd(ctx, key, member).Result()
	if err != nil {
		err = errors.ErrorCodeResponse[errors.ErrKvSever]
	}
	return err
}

func (store *RedisStore) SMembers(ctx context.Context, key string) ([][]byte, error) {
	val, err := store.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, errors.ErrorCodeResponse[errors.ErrKvSever]
	}
	if len(val) == 0 {
		return [][]byte{}, kv.KvNotFound
	}

	ret := make([][]byte, len(val))
	for i, str := range val {
		ret[i] = []byte(str)
	}

	return ret, nil
}

func (store *RedisStore) SCard(ctx context.Context, key string) (int64, error) {
	ret, err := store.client.SCard(ctx, key).Result()
	if err != nil {
		err = errors.ErrorCodeResponse[errors.ErrKvSever]
	}
	return ret, err
}

func (store *RedisStore) SRem(ctx context.Context, key string, member []byte) (bool, error) {
	ret, err := store.client.SRem(ctx, key, member).Result()
	if err != nil {
		err = errors.ErrorCodeResponse[errors.ErrKvSever]
	}
	return ret != 0, err
}

func (store *RedisStore) SIsMember(ctx context.Context, key string, member []byte) (bool, error) {
	ret, err := store.client.SIsMember(ctx, key, member).Result()
	if err != nil {
		return false, errors.ErrorCodeResponse[errors.ErrKvSever]
	}
	return ret, err
}

// delete all members of the set
func (store *RedisStore) SDelete(ctx context.Context, key string) (bool, error) {
	cnt, err := store.SCard(ctx, key)
	if err != nil {
		return false, errors.ErrorCodeResponse[errors.ErrKvSever]
	}
	if cnt == 0 {
		return false, nil
	}

	_, err = store.client.SPopN(ctx, key, cnt).Result()
	if err != nil {
		return false, errors.ErrorCodeResponse[errors.ErrKvSever]
	}
	return true, nil
}
