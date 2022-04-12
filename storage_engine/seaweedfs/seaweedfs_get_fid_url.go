package seaweedfs

import (
	"context"
	"icesos/errors"
)

func (svr *StorageEngine) GetFidUrl(ctx context.Context, fid string) (string, error) {
	volumeId, _, err := svr.parseFid(fid)
	if err != nil {
		return "", errors.ErrorCodeResponse[errors.ErrServer]
	}

	host, err := svr.getVolumeHost(ctx, volumeId)
	if err != nil {
		return "", err
	}

	url := "http://" + host + "/" + fid
	return url, nil

}
