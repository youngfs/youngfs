package storage_engine

import (
	"context"
	"io"
)

type StorageEngine interface {
	PutObject(ctx context.Context, size uint64, reader io.Reader, compress bool, hosts ...string) (string, error)
	GetObject(ctx context.Context, fid string, writer io.Writer) error
	DeleteObject(ctx context.Context, fid string) error
	GetHosts(ctx context.Context) ([]string, error)
}
