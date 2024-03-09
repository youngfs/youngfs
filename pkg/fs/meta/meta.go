package meta

import (
	"context"
	"github.com/youngfs/youngfs/pkg/fs/bucket"
	"github.com/youngfs/youngfs/pkg/fs/entry"
	"github.com/youngfs/youngfs/pkg/fs/fullpath"
)

type Store interface {
	InsertObject(ctx context.Context, ent *entry.Entry) error
	GetObject(ctx context.Context, bkt bucket.Bucket, fp fullpath.FullPath) (*entry.Entry, error)
	DeleteObjects(ctx context.Context, bkt bucket.Bucket, fp fullpath.FullPath) error
	ListObjects(ctx context.Context, bkt bucket.Bucket, fp fullpath.FullPath, recursive bool) ([]*entry.Entry, error)
}
