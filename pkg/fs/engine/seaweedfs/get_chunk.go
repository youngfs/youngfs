package seaweedfs

import (
	"compress/gzip"
	"context"
	"fmt"
	"github.com/youngfs/youngfs/pkg/errors"
	"io"
	"net/http"
)

func (e *Engine) GetChunk(ctx context.Context, id string) (io.ReadCloser, error) {
	url, err := e.getFidUrl(id)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.ErrEngineChunk.WarpErr(err)
	}
	req.Header.Add("Accept-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.ErrEngineChunk.WarpErr(err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, errors.ErrChunkNotExist
		}
		return nil, errors.ErrEngineChunk.WithMessagef("seaweedfs get failed: %s", resp.Status)
	}

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
	default:
		reader = resp.Body
	}
	defer func() {
		_ = reader.Close()
	}()

	return reader, nil
}

func (e *Engine) getFidUrl(fid string) (string, error) {
	volumeId, _, err := e.parseFid(fid)
	if err != nil {
		return "", err
	}

	host, err := e.getVolumeEndpoints(volumeId)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("http://%s/%s", host, fid), nil
}
