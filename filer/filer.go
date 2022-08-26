package filer

import (
	"context"
	"icesfs/ec/ec_store"
	"icesfs/entry"
	"icesfs/full_path"
	"icesfs/set"
	"time"
)

type FilerStore interface {
	InsertObject(ctx context.Context, ent *entry.Entry, cover bool) error
	GetObject(ctx context.Context, set set.Set, fp full_path.FullPath) (*entry.Entry, error)
	DeleteObject(ctx context.Context, set set.Set, fp full_path.FullPath, recursive bool, mtime time.Time) error
	ListObjects(ctx context.Context, set set.Set, fp full_path.FullPath) ([]entry.ListEntry, error)
	RecoverObject(ctx context.Context, frags []ec_store.Frag) error
}
