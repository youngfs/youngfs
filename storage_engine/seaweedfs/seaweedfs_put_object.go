package seaweedfs

import (
	"context"
	"icesos/command/vars"
	"icesos/errors"
	"icesos/log"
	"io"
	"net/http"
)

type PutObjectInfo struct {
	Size uint64 `json:"size"`
	ETag string `json:"eTag"`
}

func (se *StorageEngine) PutObject(ctx context.Context, size uint64, file io.Reader, hosts ...string) (string, error) {
	info, err := se.assignObject(ctx, size, hosts...)
	if err != nil {
		log.Errorw("seaweedfs put object: assign object error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "size", size, "hosts", hosts)
		return "", err
	}

	req, err := http.NewRequest("PUT", "http://"+info.Url+"/"+info.Fid, file)
	if err != nil {
		log.Errorw("seaweedfs put object: new request put error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "request url", "http://"+info.Url+"/"+info.Fid, "request", req)
		return "", errors.GetAPIErr(errors.ErrServer)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorw("seaweedfs put object: do request put error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "request url", "http://"+info.Url+"/"+info.Fid, "request", req, "response", resp)
		return "", errors.GetAPIErr(errors.ErrSeaweedFSVolume)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusCreated {
		log.Errorw("seaweedfs put object: request error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), "request url", "http://"+info.Url+"/"+info.Fid, "http code", resp.StatusCode, "request", req, "response", resp)
		return "", errors.GetAPIErr(errors.ErrSeaweedFSVolume)
	}

	return info.Fid, nil
}
