package storage_engine

import (
	"bytes"
	"context"
	"github.com/go-playground/assert/v2"
	"icesos/command/vars"
	"icesos/util"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestDeleteObject(t *testing.T) {
	client := NewStorageEngine(vars.MasterServer)
	size := uint64(5 * 1024)
	ctx := context.Background()

	b := util.RandByte(size)

	Fid, err := client.PutObject(ctx, size, bytes.NewReader(b))

	volumeId, _ := SplitFid(Fid)

	url, err := client.GetVolumeIp(ctx, volumeId)
	assert.Equal(t, err, nil)

	resp, err := http.Get("http://" + url + "/" + Fid)
	assert.Equal(t, err, nil)

	httpBody, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, err, nil)
	assert.Equal(t, httpBody, b)
	defer func() {
		_ = resp.Body.Close()
	}()

	err = client.DeleteObject(ctx, Fid)
	assert.Equal(t, err, nil)

	time.Sleep(3 * time.Second)

	resp, err = http.Get("http://" + url + "/" + Fid)
	assert.Equal(t, err, nil)
	defer func() {
		_ = resp.Body.Close()
	}()

	httpBody, err = ioutil.ReadAll(resp.Body)
	assert.Equal(t, err, nil)
	assert.Equal(t, httpBody, []byte{})
}
