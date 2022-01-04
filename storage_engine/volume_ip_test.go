package storage_engine

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

// test master: 10.249.177.55:9333
func TestGetVolumeIp(t *testing.T) {
	assert.Equal(t, volumeIpMap[1], "")
	url, err := GetVolumeIp(1)
	assert.Equal(t, url, "10.249.177.55:9222")
	assert.Equal(t, volumeIpMap[1], "10.249.177.55:9222")
	assert.Equal(t, err, nil)

	assert.Equal(t, volumeIpMap[13], "")
	url, err = GetVolumeIp(13)
	assert.Equal(t, url, "10.249.178.191:9222")
	assert.Equal(t, volumeIpMap[13], "10.249.178.191:9222")
	assert.Equal(t, err, nil)

	assert.Equal(t, volumeIpMap[17], "")
	url, err = GetVolumeIp(17)
	assert.Equal(t, url, "10.249.181.72:9222")
	assert.Equal(t, volumeIpMap[17], "10.249.181.72:9222")
	assert.Equal(t, err, nil)

	url, err = GetVolumeIp(1)
	assert.Equal(t, url, "10.249.177.55:9222")
	assert.Equal(t, err, nil)

	url, err = GetVolumeIp(100)
	assert.Equal(t, url, "")
	assert.Equal(t, err, nil)
}
