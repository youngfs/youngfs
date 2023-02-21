package filer

import (
	"context"
	"time"
	"youngfs/fs/bucket"
	"youngfs/fs/entry"
	"youngfs/fs/fullpath"
)

type FilerStore interface {
	InsertObject(ctx context.Context, ent *entry.Entry, cover bool) error
	GetObject(ctx context.Context, bkt bucket.Bucket, fp fullpath.FullPath) (*entry.Entry, error)
	DeleteObject(ctx context.Context, bkt bucket.Bucket, fp fullpath.FullPath, recursive bool, mtime time.Time) error
	ListObjects(ctx context.Context, bkt bucket.Bucket, fp fullpath.FullPath) ([]entry.ListEntry, error)
}
