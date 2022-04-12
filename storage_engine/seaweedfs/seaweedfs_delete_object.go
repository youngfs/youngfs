package seaweedfs

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"icesos/errors"
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
				volumeId, fid, err := se.parseFid(fid)
				if err != nil {
					//todo: add log
					continue
				}
				_ = se.deleteActualObject(context.Background(), volumeId, fid)
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
		return errors.ErrorCodeResponse[errors.ErrSeaweedFSVolume]
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.ErrorCodeResponse[errors.ErrSeaweedFSVolume]
	}
	if resp.StatusCode != http.StatusAccepted {
		return errors.ErrorCodeResponse[errors.ErrSeaweedFSVolume]
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	httpBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.ErrorCodeResponse[errors.ErrSeaweedFSVolume]
	}

	info := &deleteObjectInfo{}
	err = jsoniter.Unmarshal(httpBody, info)
	if err != nil {
		return errors.ErrorCodeResponse[errors.ErrSeaweedFSVolume]
	}

	return nil
}
