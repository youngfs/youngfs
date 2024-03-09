package seaweedfs

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/youngfs/youngfs/pkg/errors"
	"io"
	"net/http"
)

type deleteChunkInfo struct {
	Size uint64 `json:"size"`
}

func (e *Engine) DeleteChunk(ctx context.Context, id string) error {
	volumeId, fid, err := e.parseFid(id)
	if err != nil {
		return err
	}

	volumeEndpoint, err := e.getVolumeEndpoints(volumeId)
	if err != nil || volumeEndpoint == "" {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("http://%s/%d,%s", volumeEndpoint, volumeId, fid), nil)
	if err != nil {
		return errors.ErrEngineChunk.WarpErr(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.ErrEngineChunk.WarpErr(err)
	}
	if resp.StatusCode != http.StatusAccepted {
		return errors.ErrEngineChunk.WithMessagef("seaweedfs delete failed: %s", resp.Status)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	httpBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.ErrEngineChunk.WarpErr(err)
	}

	info := &deleteChunkInfo{}
	err = jsoniter.Unmarshal(httpBody, info)
	if err != nil {
		return errors.ErrEngineChunk.WarpErr(err)
	}

	return nil
}
