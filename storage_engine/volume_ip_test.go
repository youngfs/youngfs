package storage_engine

import (
	"context"
	"github.com/go-playground/assert/v2"
	"icesos/command/vars"
	"icesos/errors"
	"testing"
)

// test master: 10.249.177.55:9333
func TestStorageEngine_GetVolumeIp(t *testing.T) {
	client := NewStorageEngine(vars.MasterServer)
	ctx := context.Background()

	assert.Equal(t, client.volumeIpMap[1], "")
	url, err := client.GetVolumeIp(ctx, 1)
	assert.Equal(t, url, "10.249.177.55:9222")
	assert.Equal(t, client.volumeIpMap[1], "10.249.177.55:9222")
	assert.Equal(t, err, nil)

	assert.Equal(t, client.volumeIpMap[13], "")
	url, err = client.GetVolumeIp(ctx, 13)
	assert.Equal(t, url, "10.249.178.191:9222")
	assert.Equal(t, client.volumeIpMap[13], "10.249.178.191:9222")
	assert.Equal(t, err, nil)

	assert.Equal(t, client.volumeIpMap[17], "")
	url, err = client.GetVolumeIp(ctx, 17)
	assert.Equal(t, url, "10.249.181.72:9222")
	assert.Equal(t, client.volumeIpMap[17], "10.249.181.72:9222")
	assert.Equal(t, err, nil)

	url, err = client.GetVolumeIp(ctx, 1)
	assert.Equal(t, url, "10.249.177.55:9222")
	assert.Equal(t, err, nil)

	url, err = client.GetVolumeIp(ctx, 100)
	assert.Equal(t, url, "")
	assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrSeaweedFSMaster])
}
