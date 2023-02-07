package vfs

import (
	"context"
	"go.uber.org/multierr"
	"youngfs/errors"
	"youngfs/fs/entry"
	"youngfs/fs/full_path"
	"youngfs/fs/set"
)

func (vfs *VFS) insertEntry(ctx context.Context, ent *entry.Entry) error {
	b, err := ent.EncodeProto()
	if err != nil {
		return err
	}

	return vfs.kvStore.KvPut(ctx, ent.Key(), b)
}

func (vfs *VFS) getEntry(ctx context.Context, set set.Set, fp full_path.FullPath) (*entry.Entry, error) {
	key := entry.EntryKey(set, fp)

	b, err := vfs.kvStore.KvGet(ctx, key)
	if err != nil {
		return nil, err
	}

	return entry.DecodeEntryProto(b)
}

func (vfs *VFS) deleteEntry(ctx context.Context, set set.Set, fp full_path.FullPath) error {
	key := entry.EntryKey(set, fp)

	ent, err := vfs.getEntry(ctx, set, fp)
	if err != nil {
		if errors.IsKvNotFound(err) {
			return nil
		}
		return err
	}

	_, err = vfs.kvStore.KvDelete(ctx, key)
	if err != nil {
		return err
	}

	if ent.IsFile() {
		var merr error
		for _, chunk := range ent.Chunks {
			for _, frag := range chunk.Frags {
				err := vfs.storageEngine.DeleteObject(ctx, frag.Fid)
				if err != nil {
					merr = multierr.Append(merr, err)
				}
			}
		}
		if merr != nil {
			return merr
		}
	}

	return nil
}
