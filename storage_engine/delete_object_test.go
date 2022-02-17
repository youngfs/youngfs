package storage_engine

import (
	"bytes"
	"github.com/go-playground/assert/v2"
	"icesos/util"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestDeleteObject(t *testing.T) {
	size := uint64(1024 * 1024)

	b := util.RandByte(size)

	Fid, err := PutObject(size, bytes.NewReader(b))

	volumeId, fid := SplitFid(Fid)

	url, err := GetVolumeIp(volumeId)
	assert.Equal(t, err, nil)

	resp, err := http.Get("http://" + url + "/" + Fid)
	assert.Equal(t, err, nil)

	httpBody, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, err, nil)
	assert.Equal(t, httpBody, b)

	err = DeleteObject(volumeId, fid)
	assert.Equal(t, err, nil)

	resp, err = http.Get("http://" + url + "/" + Fid)
	assert.Equal(t, err, nil)

	httpBody, err = ioutil.ReadAll(resp.Body)
	assert.Equal(t, err, nil)
	assert.Equal(t, httpBody, []byte{})
}
