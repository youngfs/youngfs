package seaweedfs

import (
	"context"
	"icesos/errors"
)

func (se *StorageEngine) GetFidUrl(ctx context.Context, fid string) (string, error) {
	volumeId, _, err := se.parseFid(fid)
	if err != nil {
		return "", errors.ErrorCodeResponse[errors.ErrServer]
	}

	host, err := se.getVolumeHost(ctx, volumeId)
	if err != nil {
		return "", err
	}

	url := "http://" + host + "/" + fid
	return url, nil

}
