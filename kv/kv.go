package kv

import (
	"context"
	"github.com/go-redsync/redsync/v4"
	"icesos/errors"
	"net/http"
)

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
	ZAdd(ctx context.Context, key, member string) error
	ZCard(ctx context.Context, key string) (int64, error)
	ZRem(ctx context.Context, key, member string) (bool, error)
	ZRangeByLex(ctx context.Context, key, min, max string) ([]string, error)
	ZRemRangeByLex(ctx context.Context, key, min, max string) (bool, error)
	Incr(ctx context.Context, key string) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)
	GetNum(ctx context.Context, key string) (int64, error)
	SetNum(ctx context.Context, key string, num int64) error
	ClrNum(ctx context.Context, key string) (bool, error)
}

type KvStoreWithRedisMutex interface {
	NewMutex(name string, options ...redsync.Option) *redsync.Mutex
	KvPut(ctx context.Context, key string, val []byte) error
	KvGet(ctx context.Context, key string) ([]byte, error)
	KvDelete(ctx context.Context, key string) (bool, error)
	SAdd(ctx context.Context, key string, member []byte) error
	SMembers(ctx context.Context, key string) ([][]byte, error)
	SCard(ctx context.Context, key string) (int64, error)
	SRem(ctx context.Context, key string, member []byte) (bool, error)
	SIsMember(ctx context.Context, key string, member []byte) (bool, error)
	SDelete(ctx context.Context, key string) (bool, error)
	ZAdd(ctx context.Context, key, member string) error
	ZCard(ctx context.Context, key string) (int64, error)
	ZRem(ctx context.Context, key, member string) (bool, error)
	ZRangeByLex(ctx context.Context, key, min, max string) ([]string, error)
	ZRemRangeByLex(ctx context.Context, key, min, max string) (bool, error)
	Incr(ctx context.Context, key string) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)
	GetNum(ctx context.Context, key string) (int64, error)
	SetNum(ctx context.Context, key string, num int64) error
	ClrNum(ctx context.Context, key string) (bool, error)
}

var NotFound = errors.APIError{
	ErrorCode:      errors.ErrKvNotFound,
	Description:    "Kv not found",
	HTTPStatusCode: http.StatusContinue,
}
