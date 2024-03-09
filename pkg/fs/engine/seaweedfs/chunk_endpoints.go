package seaweedfs

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/youngfs/youngfs/pkg/errors"
	"io"
	"net/http"
)

type chunkUrl struct {
	Url       string `json:"url"`
	PublicUrl string `json:"publicUrl"`
}

type chunkEndpointsInfo struct {
	VolumeOrFileId string     `json:"volumeOrFileId"`
	Locations      []chunkUrl `json:"locations"`
	Error          string     `json:"error"`
}

func (e *Engine) getVolumeEndpoints(volumeId uint64) (string, error) {
	if val, ok := e.volumeIpMap.Load(volumeId); ok {
		if ip, ok := val.(string); ok {
			return ip, nil
		}
	}

	resp, err := http.Get(fmt.Sprintf("http://%s/dir/lookup?volumeId=%d", e.masterEndpoint, volumeId))
	if err != nil {
		return "", errors.ErrEngineChunk.WarpErr(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	httpBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.ErrEngineChunk.WarpErr(err)
	}

	info := &chunkEndpointsInfo{}
	err = jsoniter.Unmarshal(httpBody, info)
	if err != nil {
		return "", errors.ErrEngineChunk.WarpErr(err)
	}

	if info.Error != "" || len(info.Locations) != 1 {
		return "", errors.ErrEngineChunk.WithMessage("seaweedfs get volume endpoints failed")
	}

	e.volumeIpMap.Store(volumeId, info.Locations[0].Url)
	return info.Locations[0].Url, nil
}
