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

	resp1, err := http.Get("http://" + url + "/" + Fid)
	assert.Equal(t, err, nil)
	defer func() {
		_ = resp1.Body.Close()
	}()

	httpBody, err := ioutil.ReadAll(resp1.Body)
	assert.Equal(t, err, nil)
	assert.Equal(t, httpBody, b)

	err = client.DeleteObject(ctx, Fid)
	assert.Equal(t, err, nil)

	time.Sleep(3 * time.Second)

	resp2, err := http.Get("http://" + url + "/" + Fid)
	assert.Equal(t, err, nil)
	defer func() {
		_ = resp2.Body.Close()
	}()

	httpBody, err = ioutil.ReadAll(resp2.Body)
	assert.Equal(t, err, nil)
	assert.Equal(t, httpBody, []byte{})
}
