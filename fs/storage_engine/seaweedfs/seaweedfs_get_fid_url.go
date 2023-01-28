package seaweedfs

import (
	"context"
	"net/http"
	"youngfs/errors"
)

func (se *StorageEngine) GetFidUrl(ctx context.Context, fid string) (string, error) {
	link, err := se.kvStore.GetNum(ctx, fidLinkKey(fid))
	if err != nil {
		return "", errors.ErrObjectNotExist
	}

	if link <= 0 {
		return "", errors.ErrServer
	}

	volumeId, _, err := se.parseFid(ctx, fid)
	if err != nil {
		return "", err
	}

	host, err := se.getVolumeHost(ctx, volumeId)
	if err != nil {
		return "", err
	}

	url := "http://" + host + "/" + fid
	resp, err := http.Head(url)
	if err != nil {
		return "", errors.ErrObjectNotExist.WithMessage("seaweedfs get fid url: head error")
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.ErrObjectNotExist.WithMessage("seaweedfs get fid url: status code error")
	}
	return url, nil

}
