package vfs

import (
	"context"
	"icesfs/entry"
	"icesfs/full_path"
	"icesfs/kv"
	"icesfs/set"
)

func (vfs *VFS) insertEntry(ctx context.Context, ent *entry.Entry) error {
	b, err := ent.EncodeProto(ctx)
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

	return entry.DecodeEntryProto(ctx, b)
}

func (vfs *VFS) deleteEntry(ctx context.Context, set set.Set, fp full_path.FullPath) error {
	key := entry.EntryKey(set, fp)

	ent, err := vfs.getEntry(ctx, set, fp)
	if err != nil {
		if err == kv.NotFound {
			return nil
		}
		return err
	}

	_, err = vfs.kvStore.KvDelete(ctx, key)
	if err != nil {
		return err
	}

	if ent.IsFile() {
		err := vfs.storageEngine.DeleteObject(ctx, ent.Fid)
		if err != nil {
			return err
		}
		err = vfs.ecServer.DeleteObject(ctx, ent.ECid)
		if err != nil {
			return err
		}
	}

	return nil
}
