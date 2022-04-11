package storage_engine

import (
	"context"
	"github.com/go-playground/assert/v2"
	"icesos/command/vars"
	"icesos/errors"
	"testing"
)

func TestStorageEngine_AssignObject(t *testing.T) {
	client := NewStorageEngine(vars.MasterServer)
	info, err := client.AssignObject(context.Background(), 5*1024)
	assert.Equal(t, err, nil)
	assert.Equal(t, info.Url, info.PublicUrl)
	assert.Equal(t, info.Count, int64(1))
}

func TestStorageEngine_ParseFid(t *testing.T) {
	volumeId, fid, err := ParseFid("3,3fd41bd1da80")
	assert.Equal(t, volumeId, uint64(3))
	assert.Equal(t, fid, "3fd41bd1da80")
	assert.Equal(t, err, nil)

	volumeId, fid, err = ParseFid("3,3fd41bd1da80,3")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrParseFid])

	volumeId, fid, err = ParseFid("3fd41bd1da80")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrParseFid])

	volumeId, fid, err = ParseFid("")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrParseFid])

	volumeId, fid, err = ParseFid("-3,3fd41bd1da80")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrParseFid])

	volumeId, fid, err = ParseFid("3fd41bd1da80,3")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrParseFid])

	volumeId, fid, err = ParseFid("3fd41bd1da80.3")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrParseFid])
}
