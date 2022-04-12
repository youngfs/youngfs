package seaweedfs

import (
	"context"
	"icesos/errors"
)

func (svr *StorageEngine) GetFidHost(ctx context.Context, fid string) (string, error) {
	volumeId, _, err := svr.parseFid(fid)
	if err != nil {
		return "", errors.ErrorCodeResponse[errors.ErrServer]
	}

	url, err := svr.getVolumeHost(ctx, volumeId)
	if err != nil {
		return "", err
	}
	return url, nil

}
