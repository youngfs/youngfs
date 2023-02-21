package seaweedfs

import (
	"context"
	"io"
	"net/http"
	"youngfs/errors"
)

type PutObjectInfo struct {
	Size uint64 `json:"size"`
	ETag string `json:"eTag"`
}

func (se *StorageEngine) PutObject(ctx context.Context, size uint64, reader io.Reader, hosts ...string) (string, error) {
	info, err := se.assignObject(size, hosts...)
	if err != nil {
		return "", err
	}

	buffer := se.bufferPool.Get()
	defer se.bufferPool.Put(buffer)
	_, err = se.gzipWriterPool.GzipStream(buffer, reader)
	if err != nil {
		return "", errors.ErrServer.Wrap("seaweedfs put object: gzip copy")
	}

	req, err := http.NewRequest("PUT", "http://"+info.Url+"/"+info.Fid, buffer)
	if err != nil {
		return "", errors.ErrServer.Wrap("seaweedfs put object: new request put error")
	}
	req.Header.Set("Content-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.ErrSeaweedFSVolume.Wrap("seaweedfs put object: do request put error")
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	return info.Fid, nil
}
