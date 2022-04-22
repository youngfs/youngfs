package seaweedfs

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"icesos/command/vars"
	"icesos/errors"
	"icesos/log"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type deleteObjectInfo struct {
	Size uint64 `json:"size"`
}

func (se *StorageEngine) DeleteObject(ctx context.Context, fid string) error {
	se.deletionQueue.EnQueue(fid)
	return nil
}

func (se *StorageEngine) loopProcessingDeletion() {
	var deleteCnt int
	for {
		deleteCnt = 0
		se.deletionQueue.Consume(func(fids []string) {
			for _, fid := range fids {
				ctx := context.Background()

				volumeId, fid, err := se.parseFid(ctx, fid)
				if err != nil {
					log.Errorw("seaweedfs delete object: parse fid error", vars.ErrorKey, err.Error(), "fid", fid)
					continue
				}
				err = se.deleteActualObject(ctx, volumeId, fid)
				if err != nil {
					log.Errorw("seaweedfs delete object: delete actual object error", vars.ErrorKey, err.Error(), "fid", fid)
					continue
				}
				deleteCnt++
			}
		})
		if deleteCnt == 0 {
			time.Sleep(1234 * time.Millisecond)
		}
	}
}

func (se *StorageEngine) deleteActualObject(ctx context.Context, volumeId uint64, fid string) error {
	volumeIp, err := se.getVolumeHost(ctx, volumeId)
	if err != nil || volumeIp == "" {
		return err
	}

	req, err := http.NewRequest("DELETE", "http://"+volumeIp+"/"+strconv.FormatUint(volumeId, 10)+","+fid, nil)
	if err != nil {
		log.Errorw("seaweedfs delete actual object: new request delete error", vars.ErrorKey, err.Error(), "request url", "http://"+volumeIp+"/"+strconv.FormatUint(volumeId, 10)+","+fid)
		return errors.GetAPIErr(errors.ErrSeaweedFSVolume)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorw("seaweedfs delete actual object: do request delete error", vars.ErrorKey, err.Error(), "request url", "http://"+volumeIp+"/"+strconv.FormatUint(volumeId, 10)+","+fid, "request", req, "response", resp)
		return errors.GetAPIErr(errors.ErrSeaweedFSVolume)
	}
	if resp.StatusCode != http.StatusAccepted {
		log.Errorw("seaweedfs delete actual object: request error", "request url", "http://"+volumeIp+"/"+strconv.FormatUint(volumeId, 10)+","+fid, "http code", resp.StatusCode, "request", req, "response", resp)
		return errors.GetAPIErr(errors.ErrSeaweedFSVolume)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	httpBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorw("seaweedfs delete actual object: get http body error", vars.ErrorKey, err.Error(), "request url", "http://"+volumeIp+"/"+strconv.FormatUint(volumeId, 10)+","+fid, "response", resp)
		return errors.GetAPIErr(errors.ErrSeaweedFSVolume)
	}

	info := &deleteObjectInfo{}
	err = jsoniter.Unmarshal(httpBody, info)
	if err != nil {
		log.Errorw("seaweedfs delete actual object: get http body error", vars.ErrorKey, err.Error(), "request url", "http://"+volumeIp+"/"+strconv.FormatUint(volumeId, 10)+","+fid, "http body", httpBody)
		return errors.GetAPIErr(errors.ErrSeaweedFSVolume)
	}

	return nil
}
