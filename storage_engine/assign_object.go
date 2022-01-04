package storage_engine

import (
	jsoniter "github.com/json-iterator/go"
	"icesos/command/vars"
	"icesos/errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type AssignObjectInfo struct {
	Fid       string `json:"fid"` //Fid = VolumeId,Fid
	Url       string `json:"url"`
	PublicUrl string `json:"publicUrl"`
	Count     int64  `json:"count"`
}

func AssignObject(size uint64) (*AssignObjectInfo, error) {
	resp, err := http.Get("http://" + vars.MasterServer + "/dir/assign?preallocate=" + strconv.FormatUint(size, 10))
	if err != nil {
		return nil, errors.ErrorCodeResponse[errors.ErrSeaweedFSMaster]
	}

	httpBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.ErrorCodeResponse[errors.ErrSeaweedFSMaster]
	}

	assignFileInfo := &AssignObjectInfo{}
	err = jsoniter.Unmarshal(httpBody, assignFileInfo)
	if err != nil {
		return nil, errors.ErrorCodeResponse[errors.ErrSeaweedFSMaster]
	}

	return assignFileInfo, nil
}

func SplitFid(fullFid string) (uint64, string) {
	ret := strings.Split(fullFid, ",")
	if len(ret) != 2 {
		return 0, ""
	}
	volumeId, err := strconv.ParseUint(ret[0], 10, 64)
	if err != nil {
		return 0, ""
	}
	return volumeId, ret[1]
}
