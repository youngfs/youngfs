package filer

import (
	"context"
	"icesos/ec/ec_store"
	"icesos/entry"
	"icesos/full_path"
	"icesos/set"
	"time"
)

type FilerStore interface {
	InsertObject(ctx context.Context, ent *entry.Entry, cover bool) error
	GetObject(ctx context.Context, set set.Set, fp full_path.FullPath) (*entry.Entry, error)
	DeleteObject(ctx context.Context, set set.Set, fp full_path.FullPath, recursive bool, mtime time.Time) error
	ListObjects(ctx context.Context, set set.Set, fp full_path.FullPath) ([]entry.ListEntry, error)
	RecoverObject(ctx context.Context, frags []ec_store.Frag) error
}
