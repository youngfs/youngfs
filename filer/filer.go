package filer

import (
	"context"
	"icesos/entry"
	"icesos/full_path"
	"icesos/set"
)

type FilerStore interface {
	InsertObject(ctx context.Context, ent *entry.Entry, cover bool) error
	GetObject(ctx context.Context, set set.Set, fp full_path.FullPath) (*entry.Entry, error)
	DeleteObject(ctx context.Context, set set.Set, fp full_path.FullPath, recursive bool) error
	ListObjects(ctx context.Context, set set.Set, fp full_path.FullPath) ([]entry.Entry, error)
}
