package seaweedfs

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/youngfs/youngfs/errors"
	"github.com/youngfs/youngfs/util"
	"io"
	"net/http"
	"strconv"
	"sync/atomic"
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
	var deleteCnt int64
	lce := util.NewLimitedConcurrentExecutor(8)
	for {
		deleteCnt = 0
		se.deletionQueue.Consume(func(fids []string) {
			for _, id := range fids {
				fid := id
				lce.Execute(func() {
					err := se.deleteActualObject(fid)
					if err != nil {
						//log.Errorw("seaweedfs delete object: delete actual object error", vars.ErrorKey, err.Error(), "fid", fid)
						return
					}

					atomic.AddInt64(&deleteCnt, 1)
				})

			}
			lce.Wait()
		})
		if deleteCnt == 0 {
			time.Sleep(time.Second)
		}
	}
}

func (se *StorageEngine) deleteActualObject(fullFid string) error {
	volumeId, fid, err := se.parseFid(fullFid)
	if err != nil {
		return err
	}

	volumeIp, err := se.getVolumeHost(volumeId)
	if err != nil || volumeIp == "" {
		return err
	}

	req, err := http.NewRequest("DELETE", "http://"+volumeIp+"/"+strconv.FormatUint(volumeId, 10)+","+fid, nil)
	if err != nil {
		return errors.ErrSeaweedFSVolume.Wrap("seaweedfs delete actual object: new request delete error")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.ErrSeaweedFSVolume.Wrap("seaweedfs delete actual object: do request delete error")
	}
	if resp.StatusCode != http.StatusAccepted {
		return errors.ErrSeaweedFSVolume.Wrap("seaweedfs delete actual object: request error")
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	httpBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.ErrSeaweedFSVolume.Wrap("seaweedfs delete actual object: get http body error")
	}

	info := &deleteObjectInfo{}
	err = jsoniter.Unmarshal(httpBody, info)
	if err != nil {
		return errors.ErrSeaweedFSVolume.Wrap("seaweedfs delete actual object: get http body error")
	}

	return nil
}
