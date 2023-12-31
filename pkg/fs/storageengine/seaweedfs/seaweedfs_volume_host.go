package seaweedfs

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/youngfs/youngfs/pkg/errors"
	"io"
	"net/http"
	"strconv"
)

type volumeUrl struct {
	Url       string `json:"url"`
	PublicUrl string `json:"publicUrl"`
}

type volumeIpInfo struct {
	VolumeOrFileId string      `json:"volumeOrFileId"`
	Locations      []volumeUrl `json:"locations"`
	Error          string      `json:"error"`
}

func (se *StorageEngine) getVolumeHost(volumeId uint64) (string, error) {
	if val, ok := se.volumeIpMap.Load(volumeId); ok {
		if ip, ok := val.(string); ok {
			return ip, nil
		}
	}

	resp, err := http.Get("http://" + se.masterServer + "/dir/lookup?volumeId=" + strconv.FormatUint(volumeId, 10))
	if err != nil {
		return "", errors.ErrSeaweedFSMaster.Wrap("seaweedfs get volume host : http get error")
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	httpBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.ErrSeaweedFSMaster.Wrap("seaweedfs get volume host : get http body error")
	}

	info := &volumeIpInfo{}
	err = jsoniter.Unmarshal(httpBody, info)
	if err != nil {
		return "", errors.ErrSeaweedFSMaster.Wrap("seaweedfs get volume host : http body unmarshal error")
	}

	if info.Error != "" || len(info.Locations) != 1 {
		return "", errors.ErrSeaweedFSMaster.Wrap("seaweedfs get volume host : http body unmarshal error")
	}

	se.volumeIpMap.Store(volumeId, info.Locations[0].Url)
	return info.Locations[0].Url, nil
}
