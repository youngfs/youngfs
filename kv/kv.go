package kv

import (
	"context"
	"github.com/go-redsync/redsync/v4"
)

type KvStore interface {
	KvPut(ctx context.Context, key string, val []byte) error
	KvGet(ctx context.Context, key string) ([]byte, error)
	KvDelete(ctx context.Context, key string) (bool, error)
}

type KvSetStore interface {
	KvStore
	ZAdd(ctx context.Context, key, member string) error
	ZCard(ctx context.Context, key string) (int64, error)
	ZRem(ctx context.Context, key, member string) (bool, error)
	ZRangeByLex(ctx context.Context, key, min, max string) ([]string, error)
	ZRemRangeByLex(ctx context.Context, key, min, max string) (bool, error)
	ZIsMember(ctx context.Context, key, member string) (bool, error)
}

type KvSetStoreWithRedisMutex interface {
	KvSetStore
	NewMutex(name string, options ...redsync.Option) *redsync.Mutex
}
