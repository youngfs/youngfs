package storage_engine

import (
	jsoniter "github.com/json-iterator/go"
	"icesos/command/vars"
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
	VolumeOrFileId string       `json:"volumeOrFileId"`
	Locations      []*volumeUrl `json:"locations"`
	Error          string       `json:"error"`
}

var volumeIpMap = map[uint64]string{}

func GetVolumeIp(volumeId uint64) (string, error) {
	if volumeIpMap[volumeId] != "" {
		return volumeIpMap[volumeId], nil
	}

	resp, err := http.Get("http://" + vars.MasterServer + "/dir/lookup?volumeId=" + strconv.FormatUint(volumeId, 10))
	if err != nil {
		return "", errors.ErrorCodeResponse[errors.ErrSeaweedFSMaster]
	}

	httpBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.ErrorCodeResponse[errors.ErrSeaweedFSMaster]
	}

	info := &volumeIpInfo{}
	err = jsoniter.Unmarshal(httpBody, info)

	if info.Error != "" || len(info.Locations) != 1 {
		return "", errors.ErrorCodeResponse[errors.ErrServer]
	}

	volumeIpMap[volumeId] = info.Locations[0].Url
	return volumeIpMap[volumeId], nil
}
