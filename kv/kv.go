package kv

import "context"

type KvStore interface {
	KvPut(ctx context.Context, key string, val []byte) error
	KvGet(ctx context.Context, key string) ([]byte, error)
	KvDelete(ctx context.Context, key string) (bool, error)
	SAdd(ctx context.Context, key string, member []byte) error
	SMembers(ctx context.Context, key string) ([][]byte, error)
	SCard(ctx context.Context, key string) (int64, error)
	SRem(ctx context.Context, key string, member []byte) (bool, error)
	SIsMember(ctx context.Context, key string, member []byte) (bool, error)
	SDelete(ctx context.Context, key string) (bool, error)
}
