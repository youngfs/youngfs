package seaweedfs

import (
	"context"
	"icesos/command/vars"
	"icesos/errors"
	"icesos/log"
)

func (se *StorageEngine) GetFidUrl(ctx context.Context, fid string) (string, error) {
	volumeId, _, err := se.parseFid(ctx, fid)
	if err != nil {
		log.Errorw("seaweedfs get fid url: parse fid error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "fid", fid)
		return "", errors.GetAPIErr(errors.ErrServer)
	}

	host, err := se.getVolumeHost(ctx, volumeId)
	if err != nil {
		return "", err
	}

	url := "http://" + host + "/" + fid
	return url, nil

}
