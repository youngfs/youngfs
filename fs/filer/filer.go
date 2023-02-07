package filer

import (
	"context"
	"time"
	"youngfs/fs/entry"
	"youngfs/fs/full_path"
	"youngfs/fs/set"
)

type FilerStore interface {
	InsertObject(ctx context.Context, ent *entry.Entry, cover bool) error
	GetObject(ctx context.Context, set set.Set, fp full_path.FullPath) (*entry.Entry, error)
	DeleteObject(ctx context.Context, set set.Set, fp full_path.FullPath, recursive bool, mtime time.Time) error
	ListObjects(ctx context.Context, set set.Set, fp full_path.FullPath) ([]entry.ListEntry, error)
}
