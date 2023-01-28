package storage_engine

import (
	"context"
	"io"
)

type StorageEngine interface {
	PutObject(ctx context.Context, size uint64, file io.Reader, fileName string, compress bool, hosts ...string) (string, error)
	DeleteObject(ctx context.Context, fid string) error
	GetFidUrl(ctx context.Context, fid string) (string, error)
	GetHosts(ctx context.Context) ([]string, error)
	AddLink(ctx context.Context, fid string) error
}
