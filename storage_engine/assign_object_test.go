package storage_engine

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestAssignObject(t *testing.T) {
	info, err := AssignObject(1024 * 1024)
	assert.Equal(t, err, nil)
	assert.Equal(t, info.Url, info.PublicUrl)
	assert.Equal(t, info.Count, int64(1))
}

func TestSplitFid(t *testing.T) {
	volumeId, fid := SplitFid("3,3fd41bd1da80")
	assert.Equal(t, volumeId, uint64(3))
	assert.Equal(t, fid, "3fd41bd1da80")

	volumeId, fid = SplitFid("3,3fd41bd1da80,3")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")

	volumeId, fid = SplitFid("3fd41bd1da80")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")

	volumeId, fid = SplitFid("")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")

	volumeId, fid = SplitFid("-3,3fd41bd1da80")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")

	volumeId, fid = SplitFid("3fd41bd1da80,3")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")

	volumeId, fid = SplitFid("3fd41bd1da80.3")
	assert.Equal(t, volumeId, uint64(0))
	assert.Equal(t, fid, "")
}
