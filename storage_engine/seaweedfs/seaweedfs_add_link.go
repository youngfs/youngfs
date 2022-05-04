package seaweedfs

import (
	"context"
	"icesos/command/vars"
	"icesos/log"
)

func fidLinkKey(fid string) string {
	return fid + seaweedfsFidLinkKey
}

func (se *StorageEngine) AddLink(ctx context.Context, fid string) error {
	_, err := se.kvStore.Incr(ctx, fidLinkKey(fid))
	if err != nil {
		log.Errorw("seaweedfs add link", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), "fid", fid)
		return err
	}
	return nil
}

func (se *StorageEngine) delLink(ctx context.Context, fid string) (int64, error) {
	ret, err := se.kvStore.Decr(ctx, fidLinkKey(fid))
	if err != nil {
		log.Errorw("seaweedfs del link", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "fid", fid)
		return -1, err
	}
	return ret, nil
}
