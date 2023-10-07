package seaweedfs

import (
	"github.com/go-playground/assert/v2"
	"github.com/youngfs/youngfs/errors"
	"github.com/youngfs/youngfs/vars"
	"testing"
)

func TestSeaweedFS_assignObject(t *testing.T) {
	client := NewStorageEngine(vars.SeaweedFSMaster)
	info, err := client.assignObject(5 * 1024)
	assert.Equal(t, err, nil)
	assert.Equal(t, info.Url, info.PublicUrl)
	assert.Equal(t, info.Count, int64(1))
}

func TestSeaweedFS_parseFid(t *testing.T) {
	client := NewStorageEngine(vars.SeaweedFSMaster)
	volumeId, fid, err := client.parseFid("3,3fd41bd1da80")
	assert.Equal(t, volumeId, uint64(3))
	assert.Equal(t, fid, "3fd41bd1da80")
	assert.Equal(t, err, nil)

	volumeId, fid, err = client.parseFid("3,3fd41bd1da80,3")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, errors.Is(err, errors.ErrServer), true)

	volumeId, fid, err = client.parseFid("3fd41bd1da80")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, errors.Is(err, errors.ErrServer), true)

	volumeId, fid, err = client.parseFid("")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, errors.Is(err, errors.ErrServer), true)

	volumeId, fid, err = client.parseFid("-3,3fd41bd1da80")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, errors.Is(err, errors.ErrServer), true)

	volumeId, fid, err = client.parseFid("3fd41bd1da80,3")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, errors.Is(err, errors.ErrServer), true)

	volumeId, fid, err = client.parseFid("3fd41bd1da80.3")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
	assert.Equal(t, errors.Is(err, errors.ErrServer), true)
}
