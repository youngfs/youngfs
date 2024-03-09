package seaweedfs

import (
	"context"
	"fmt"
	"github.com/youngfs/youngfs/pkg/errors"
	"io"
	"net/http"
)

type PutChunkInfo struct {
	Size uint64 `json:"size"`
	ETag string `json:"eTag"`
}

func (e *Engine) PutChunk(ctx context.Context, reader io.Reader, endpoints ...string) (string, error) {
	info, err := e.assignChunk(ctx, endpoints...)
	if err != nil {
		return "", err
	}

	buffer := e.bufferPool.Get()
	defer e.bufferPool.Put(buffer)
	_, err = e.gzipWriterPool.GzipStream(buffer, reader)
	if err != nil {
		return "", errors.ErrEngineChunk.WarpErr(err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("http://%s/%s", info.Url, info.Fid), buffer)
	if err != nil {
		return "", errors.ErrEngineChunk.WarpErr(err)
	}
	req.Header.Set("Content-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.ErrEngineChunk.WarpErr(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	return info.Fid, nil
}
