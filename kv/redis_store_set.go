package kv

import (
	"context"
	"icesos/errors"
)

func (store *redisStore) SAdd(key string, member []byte) error {
	_, err := store.client.SAdd(context.Background(), key, member).Result()
	if err != nil {
		err = errors.ErrorCodeResponse[errors.ErrKvSever]
	}
	return err
}

func (store *redisStore) SMembers(key string) ([][]byte, error) {
	val, err := store.client.SMembers(context.Background(), key).Result()
	if err != nil {
		return nil, errors.ErrorCodeResponse[errors.ErrKvSever]
	}
	if len(val) == 0 {
		return nil, nil
	}

	ret := make([][]byte, len(val))
	for i, str := range val {
		ret[i] = []byte(str)
	}

	return ret, nil
}

func (store *redisStore) SCard(key string) (int64, error) {
	ret, err := store.client.SCard(context.Background(), key).Result()
	if err != nil {
		err = errors.ErrorCodeResponse[errors.ErrKvSever]
	}
	return ret, err
}

func (store *redisStore) SRem(key string, member []byte) (bool, error) {
	ret, err := store.client.SRem(context.Background(), key, member).Result()
	if err != nil {
		err = errors.ErrorCodeResponse[errors.ErrKvSever]
	}
	return ret != 0, err
}

func (store *redisStore) SIsMember(key string, member []byte) (bool, error) {
	ret, err := store.client.SIsMember(context.Background(), key, member).Result()
	if err != nil {
		err = errors.ErrorCodeResponse[errors.ErrKvSever]
	}
	return ret, err
}

// delete all members of the set
func (store *redisStore) SDelete(key string) (bool, error) {
	cnt, err := store.SCard(key)
	if err != nil || cnt == 0 {
		return false, err
	}

	_, err = store.client.SPopN(context.Background(), key, cnt).Result()
	if err != nil {
		return false, errors.ErrorCodeResponse[errors.ErrKvSever]
	}
	return true, nil
}
