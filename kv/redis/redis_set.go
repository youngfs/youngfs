package redis

import (
	"context"
	"youngfs/errors"
)

func (store *KvStore) SAdd(ctx context.Context, key string, member []byte) error {
	_, err := store.client.SAdd(ctx, key, member).Result()
	if err != nil {
		err = errors.ErrKvSever.Wrap("redis sadd error")
	}
	return err
}

func (store *KvStore) SMembers(ctx context.Context, key string) ([][]byte, error) {
	val, err := store.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, errors.ErrKvSever.Wrap("redis smembers error")
	}
	if len(val) == 0 {
		return [][]byte{}, errors.ErrKvNotFound
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
		err = errors.ErrKvSever.Wrap("redis scard error")
	}
	return ret, err
}

func (store *KvStore) SRem(ctx context.Context, key string, member []byte) (bool, error) {
	ret, err := store.client.SRem(ctx, key, member).Result()
	if err != nil {
		err = errors.ErrKvSever.Wrap("redis srem error")
	}
	return ret != 0, err
}

func (store *KvStore) SIsMember(ctx context.Context, key string, member []byte) (bool, error) {
	ret, err := store.client.SIsMember(ctx, key, member).Result()
	if err != nil {
		return false, errors.ErrKvSever.Wrap("redis sismember error")
	}
	return ret, err
}

// delete all members of the bucket
func (store *KvStore) SDelete(ctx context.Context, key string) (bool, error) {
	cnt, err := store.SCard(ctx, key)
	if err != nil {
		return false, errors.ErrKvSever.Wrap("redis scard error")
	}
	if cnt == 0 {
		return false, nil
	}

	_, err = store.client.SPopN(ctx, key, cnt).Result()
	if err != nil {
		return false, errors.ErrKvSever.Wrap("redis spopn error")
	}
	return true, nil
}
