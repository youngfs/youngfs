package seaweedfs

import (
	"context"
	"github.com/go-playground/assert/v2"
	"icesos/command/vars"
	"icesos/errors"
	"testing"
)

func TestSeaweedFS_assignObject(t *testing.T) {
	client := NewStorageEngine(vars.MasterServer)
	info, err := client.assignObject(context.Background(), 5*1024)
	assert.Equal(t, err, nil)
	assert.Equal(t, info.Url, info.PublicUrl)
	assert.Equal(t, info.Count, int64(1))
}

func TestSeaweedFS_parseFid(t *testing.T) {
	client := NewStorageEngine(vars.MasterServer)
	volumeId, fid, err := client.parseFid("3,3fd41bd1da80")
	assert.Equal(t, volumeId, uint64(3))
	assert.Equal(t, fid, "3fd41bd1da80")
	assert.Equal(t, err, nil)

	volumeId, fid, err = client.parseFid("3,3fd41bd1da80,3")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrParseFid])

	volumeId, fid, err = client.parseFid("3fd41bd1da80")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrParseFid])

	volumeId, fid, err = client.parseFid("")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrParseFid])

	volumeId, fid, err = client.parseFid("-3,3fd41bd1da80")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrParseFid])

	volumeId, fid, err = client.parseFid("3fd41bd1da80,3")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrParseFid])

	volumeId, fid, err = client.parseFid("3fd41bd1da80.3")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrParseFid])
}
