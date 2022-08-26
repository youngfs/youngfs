package seaweedfs

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"icesfs/command/vars"
	"icesfs/errors"
	"icesfs/log"
	"io/ioutil"
	"math/rand"
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

func (se *StorageEngine) assignObject(ctx context.Context, size uint64, hosts ...string) (*assignObjectInfo, error) {
	hostReq, host := "", ""
	if len(hosts) > 0 {
		host = hosts[rand.Intn(len(hosts))]
		hostReq = "&dataCenter=DefaultDataCenter&rack=DefaultRack&dataNode=" + host
	}

	resp, err := http.Get("http://" + se.masterServer + "/dir/assign?preallocate=" + strconv.FormatUint(size, 10) + hostReq)
	if err != nil {
		log.Errorw("seaweedfs assign object: http get error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "request url", "http://"+se.masterServer+"/dir/assign?preallocate="+strconv.FormatUint(size, 10)+hostReq, "response", resp)
		return nil, errors.GetAPIErr(errors.ErrSeaweedFSMaster)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	httpBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorw("seaweedfs assign object: get http body error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "request url", "http://"+se.masterServer+"/dir/assign?preallocate="+strconv.FormatUint(size, 10)+hostReq, "response", resp)
		return nil, errors.GetAPIErr(errors.ErrSeaweedFSMaster)
	}

	assignFileInfo := &assignObjectInfo{}
	err = jsoniter.Unmarshal(httpBody, assignFileInfo)
	if err != nil {
		log.Errorw("seaweedfs assign object: http body unmarshal error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "request url", "http://"+se.masterServer+"/dir/assign?preallocate="+strconv.FormatUint(size, 10)+hostReq, "http body", httpBody)
		return nil, errors.GetAPIErr(errors.ErrSeaweedFSMaster)
	}

	if host != "" {
		if host != assignFileInfo.Url && host != assignFileInfo.PublicUrl {
			log.Errorw("seaweedfs assign object: request host error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), "request url", "http://"+se.masterServer+"/dir/assign?preallocate="+strconv.FormatUint(size, 10)+hostReq, "request host", host, "request hosts", hosts, "assign file info", assignFileInfo)
			return nil, errors.GetAPIErr(errors.ErrSeaweedFSMaster)
		}
	}

	return assignFileInfo, nil
}

func (se *StorageEngine) parseFid(ctx context.Context, fullFid string) (uint64, string, error) {
	ret := strings.Split(fullFid, ",")
	if len(ret) != 2 {
		log.Errorw("seaweedfs parse fid error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), "fid", fullFid)
		return 0, "", errors.GetAPIErr(errors.ErrParseFid)
	}
	volumeId, err := strconv.ParseUint(ret[0], 10, 64)
	if err != nil {
		log.Errorw("seaweedfs parse fid error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "fid", fullFid)
		return 0, "", errors.GetAPIErr(errors.ErrParseFid)
	}
	return volumeId, ret[1], nil
}
