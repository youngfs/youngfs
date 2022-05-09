package seaweedfs

import (
	"context"
	"icesos/command/vars"
	"icesos/errors"
	"icesos/log"
	"net/http"
)

func (se *StorageEngine) GetFidUrl(ctx context.Context, fid string) (string, error) {
	link, err := se.kvStore.GetNum(ctx, fidLinkKey(fid))
	if err != nil {
		return "", errors.GetAPIErr(errors.ErrObjectNotExist)
	}

	if link <= 0 {
		return "", errors.GetAPIErr(errors.ErrServer)
	}

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
	resp, err := http.Head(url)
	if err != nil {
		log.Errorw("seaweedfs get fid url: head error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "fid", fid)
		return "", errors.GetAPIErr(errors.ErrObjectNotExist)
	}
	if resp.StatusCode != http.StatusOK {
		log.Errorw("seaweedfs get fid url: status code error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), "fid", fid, "http code", resp.StatusCode)
		return "", errors.GetAPIErr(errors.ErrObjectNotExist)
	}
	return url, nil

}
