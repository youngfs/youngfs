package seaweedfs

import (
	jsoniter "github.com/json-iterator/go"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"youngfs/errors"
)

type assignObjectInfo struct {
	Fid       string `json:"fid"` //Fid = VolumeId,Fid
	Url       string `json:"url"`
	PublicUrl string `json:"publicUrl"`
	Count     int64  `json:"count"`
}

func (se *StorageEngine) assignObject(size uint64, hosts ...string) (*assignObjectInfo, error) {
	hostReq, host := "", ""
	if len(hosts) > 0 {
		host = hosts[rand.Intn(len(hosts))]
		hostReq = "&dataCenter=DefaultDataCenter&rack=DefaultRack&dataNode=" + host
	}

	resp, err := http.Get("http://" + se.masterServer + "/dir/assign?preallocate=" + strconv.FormatUint(size, 10) + hostReq)
	if err != nil {
		return nil, errors.ErrSeaweedFSMaster.Wrap("seaweedfs assign object: http get error")
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	httpBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.ErrSeaweedFSMaster.Wrap("seaweedfs assign object: get http body error")
	}

	assignFileInfo := &assignObjectInfo{}
	err = jsoniter.Unmarshal(httpBody, assignFileInfo)
	if err != nil {
		return nil, errors.ErrSeaweedFSMaster.Wrap("seaweedfs assign object: http body unmarshal error")
	}

	if host != "" {
		if host != assignFileInfo.Url && host != assignFileInfo.PublicUrl {
			return nil, errors.ErrSeaweedFSMaster.Wrap("seaweedfs assign object: request host error")
		}
	}

	return assignFileInfo, nil
}

func (se *StorageEngine) parseFid(fullFid string) (uint64, string, error) {
	ret := strings.Split(fullFid, ",")
	if len(ret) != 2 {
		return 0, "", errors.ErrServer.Wrap("seaweedfs parse fid error")
	}
	volumeId, err := strconv.ParseUint(ret[0], 10, 64)
	if err != nil {
		return 0, "", errors.ErrServer.Wrap("seaweedfs parse fid error")
	}
	return volumeId, ret[1], nil
}
