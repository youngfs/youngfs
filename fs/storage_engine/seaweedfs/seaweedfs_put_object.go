package seaweedfs

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"youngfs/errors"
)

type PutObjectInfo struct {
	Size uint64 `json:"size"`
	ETag string `json:"eTag"`
}

func (se *StorageEngine) PutObject(ctx context.Context, size uint64, reader io.Reader, compress bool, hosts ...string) (string, error) {
	info, err := se.assignObject(size, hosts...)
	if err != nil {
		return "", err
	}

	if compress {
		b := &bytes.Buffer{}
		_, err := se.gzipWriterPool.GzipStream(b, reader)
		if err != nil {
			return "", errors.ErrServer.Wrap("seaweedfs put object: gzip copy")
		}
		reader = b
	}

	req, err := http.NewRequest("PUT", "http://"+info.Url+"/"+info.Fid, reader)
	if err != nil {
		return "", errors.ErrServer.Wrap("seaweedfs put object: new request put error")
	}
	if compress {
		req.Header.Set("Content-Encoding", "gzip")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.ErrSeaweedFSVolume.Wrap("seaweedfs put object: do request put error")
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	return info.Fid, nil
}
