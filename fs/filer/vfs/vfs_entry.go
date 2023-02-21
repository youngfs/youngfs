package vfs

import (
	"context"
	"go.uber.org/multierr"
	"youngfs/errors"
	"youngfs/fs/bucket"
	"youngfs/fs/entry"
	"youngfs/fs/fullpath"
)

func (vfs *VFS) insertEntry(ctx context.Context, ent *entry.Entry) error {
	b, err := ent.EncodeProto()
	if err != nil {
		return err
	}

	return vfs.kvStore.KvPut(ctx, ent.Key(), b)
}

func (vfs *VFS) getEntry(ctx context.Context, bkt bucket.Bucket, fp fullpath.FullPath) (*entry.Entry, error) {
	key := entry.EntryKey(bkt, fp)

	b, err := vfs.kvStore.KvGet(ctx, key)
	if err != nil {
		return nil, err
	}

	return entry.DecodeEntryProto(b)
}

func (vfs *VFS) deleteEntry(ctx context.Context, bkt bucket.Bucket, fp fullpath.FullPath) error {
	key := entry.EntryKey(bkt, fp)

	ent, err := vfs.getEntry(ctx, bkt, fp)
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
