package storage_engine

import (
	"context"
	"icesos/errors"
	"io"
	"net/http"
)

type PutObjectInfo struct {
	Size uint64 `json:"size"`
	ETag string `json:"eTag"`
}

func (svr *StorageEngine) PutObject(ctx context.Context, size uint64, file io.Reader) (string, error) {
	info, err := svr.AssignObject(ctx, size)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("PUT", "http://"+info.Url+"/"+info.Fid, file)
	if err != nil {
		return "", errors.ErrorCodeResponse[errors.ErrServer]
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.ErrorCodeResponse[errors.ErrSeaweedFSVolume]
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusCreated {
		return "", errors.ErrorCodeResponse[errors.ErrSeaweedFSVolume]
	}

	return info.Fid, nil
}
