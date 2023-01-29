package seaweedfs

import (
	"context"
	"youngfs/errors"
)

func fidLinkKey(fid string) string {
	return fid + seaweedfsFidLinkKey
}

func (se *StorageEngine) AddLink(ctx context.Context, fid string) error {
	_, err := se.kvStore.Incr(ctx, fidLinkKey(fid))
	if err != nil {
		return errors.Wrap(err, "seaweedfs add link")
	}
	return nil
}

func (se *StorageEngine) delLink(ctx context.Context, fid string) (int64, error) {
	ret, err := se.kvStore.Decr(ctx, fidLinkKey(fid))
	if err != nil {
		return -1, errors.Wrap(err, "seaweedfs del link")
	}
	return ret, nil
}
