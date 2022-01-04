package storage_engine

import (
	"bytes"
	"github.com/go-playground/assert/v2"
	jsoniter "github.com/json-iterator/go"
	"icesos/util"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestDeleteObject(t *testing.T) {
	size := uint64(1024 * 1024)

	info, err := AssignObject()
	assert.Equal(t, err, nil)

	b := util.RandByte(size)
	req, err := http.NewRequest("PUT", "http://"+info.Url+"/"+info.Fid, bytes.NewReader(b))
	assert.Equal(t, err, nil)

	resp, err := http.DefaultClient.Do(req)
	assert.Equal(t, err, nil)

	httpBody, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, err, nil)

	putInfo := &PutObjectInfo{}
	err = jsoniter.Unmarshal(httpBody, putInfo)
	assert.Equal(t, err, nil)
	assert.Equal(t, putInfo.Size, size)

	resp, err = http.Get("http://" + info.Url + "/" + info.Fid)
	assert.Equal(t, err, nil)

	httpBody, err = ioutil.ReadAll(resp.Body)
	assert.Equal(t, err, nil)
	assert.Equal(t, httpBody, b)

	volumeId, fid := SplitFid(info.Fid)
	err = DeleteObject(volumeId, fid, size)
	assert.Equal(t, err, nil)

	resp, err = http.Get("http://" + info.Url + "/" + info.Fid)
	assert.Equal(t, err, nil)

	httpBody, err = ioutil.ReadAll(resp.Body)
	assert.Equal(t, err, nil)
	assert.Equal(t, httpBody, []byte{})
}
