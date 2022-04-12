package seaweedfs

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"icesos/errors"
	"io/ioutil"
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

func (svr *StorageEngine) getVolumeHost(ctx context.Context, volumeId uint64) (string, error) {
	if svr.volumeIpMap[volumeId] != "" {
		return svr.volumeIpMap[volumeId], nil
	}

	resp, err := http.Get("http://" + svr.masterServer + "/dir/lookup?volumeId=" + strconv.FormatUint(volumeId, 10))
	if err != nil {
		return "", errors.ErrorCodeResponse[errors.ErrSeaweedFSMaster]
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	httpBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.ErrorCodeResponse[errors.ErrSeaweedFSMaster]
	}

	info := &volumeIpInfo{}
	err = jsoniter.Unmarshal(httpBody, info)

	if info.Error != "" || len(info.Locations) != 1 {
		return "", errors.ErrorCodeResponse[errors.ErrSeaweedFSMaster]
	}

	svr.volumeIpMap[volumeId] = info.Locations[0].Url
	return svr.volumeIpMap[volumeId], nil
}
