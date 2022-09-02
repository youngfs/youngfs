package seaweedfs

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"icesfs/command/vars"
	"icesfs/errors"
	"icesfs/log"
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

func (se *StorageEngine) getVolumeHost(ctx context.Context, volumeId uint64) (string, error) {
	if se.volumeIpMap[volumeId] != "" {
		return se.volumeIpMap[volumeId], nil
	}

	resp, err := http.Get("http://" + se.masterServer + "/dir/lookup?volumeId=" + strconv.FormatUint(volumeId, 10))
	if err != nil {
		log.Errorw("seaweedfs get volume host : http get error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "request url", "http://"+se.masterServer+"/dir/lookup?volumeId="+strconv.FormatUint(volumeId, 10), "response", resp)
		return "", errors.GetAPIErr(errors.ErrSeaweedFSMaster)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	httpBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorw("seaweedfs get volume host : get http body error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "request url", "http://"+se.masterServer+"/dir/lookup?volumeId="+strconv.FormatUint(volumeId, 10), "response", resp)
		return "", errors.GetAPIErr(errors.ErrSeaweedFSMaster)
	}

	info := &volumeIpInfo{}
	err = jsoniter.Unmarshal(httpBody, info)
	if err != nil {
		log.Errorw("seaweedfs get volume host : http body unmarshal error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "request url", "http://"+se.masterServer+"/dir/lookup?volumeId="+strconv.FormatUint(volumeId, 10))
	}

	if info.Error != "" || len(info.Locations) != 1 {
		log.Errorw("seaweedfs get volume host : server error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), "request url", "http://"+se.masterServer+"/dir/lookup?volumeId="+strconv.FormatUint(volumeId, 10), "volume ip info", info)
		return "", errors.GetAPIErr(errors.ErrSeaweedFSMaster)
	}

	se.volumeIpMap[volumeId] = info.Locations[0].Url
	return se.volumeIpMap[volumeId], nil
}
