package storage_engine

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"icesos/errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type assignObjectInfo struct {
	Fid       string `json:"fid"` //Fid = VolumeId,Fid
	Url       string `json:"url"`
	PublicUrl string `json:"publicUrl"`
	Count     int64  `json:"count"`
}

func (svr *StorageEngine) AssignObject(ctx context.Context, size uint64) (*assignObjectInfo, error) {
	resp, err := http.Get("http://" + svr.masterServer + "/dir/assign?preallocate=" + strconv.FormatUint(size, 10))
	if err != nil {
		return nil, errors.ErrorCodeResponse[errors.ErrSeaweedFSMaster]
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	httpBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.ErrorCodeResponse[errors.ErrSeaweedFSMaster]
	}

	assignFileInfo := &assignObjectInfo{}
	err = jsoniter.Unmarshal(httpBody, assignFileInfo)
	if err != nil {
		return nil, errors.ErrorCodeResponse[errors.ErrSeaweedFSMaster]
	}

	return assignFileInfo, nil
}

func ParseFid(fullFid string) (uint64, string, error) {
	ret := strings.Split(fullFid, ",")
	if len(ret) != 2 {
		return 0, "", errors.ErrorCodeResponse[errors.ErrParseFid]
	}
	volumeId, err := strconv.ParseUint(ret[0], 10, 64)
	if err != nil {
		return 0, "", errors.ErrorCodeResponse[errors.ErrParseFid]
	}
	return volumeId, ret[1], nil
}
