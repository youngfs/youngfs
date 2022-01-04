package storage_engine

import (
	jsoniter "github.com/json-iterator/go"
	"icesos/errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

type deleteObjectInfo struct {
	Size uint64 `json:"size"`
}

func DeleteObject(volumeId uint64, fid string, size uint64) error {
	volumeIp, err := GetVolumeIp(volumeId)
	if err != nil || volumeIp == "" {
		return err
	}

	req, err := http.NewRequest("DELETE", "http://"+volumeIp+"/"+strconv.FormatUint(volumeId, 10)+","+fid, nil)
	if err != nil {
		return errors.ErrorCodeResponse[errors.ErrSeaweedFSVolume]
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.ErrorCodeResponse[errors.ErrSeaweedFSVolume]
	}

	httpBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.ErrorCodeResponse[errors.ErrSeaweedFSVolume]
	}

	info := &deleteObjectInfo{}
	err = jsoniter.Unmarshal(httpBody, info)

	if info.Size != size+deleteOffset {
		return errors.ErrorCodeResponse[errors.ErrSeaweedFSVolume]
	}

	return nil
}
