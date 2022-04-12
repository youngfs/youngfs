package seaweedfs

import (
	"context"
	"github.com/go-playground/assert/v2"
	"icesos/command/vars"
	"icesos/errors"
	"testing"
)

func TestSeaweedFS_GetVolumeHost(t *testing.T) {
	client := NewStorageEngine(vars.MasterServer)
	ctx := context.Background()

	hosts, err := client.GetHosts(ctx)
	assert.Equal(t, err, nil)

	volumes := make(map[uint64]bool)
	hostSet := make(map[string]bool)
	for _, u := range hosts {
		hostSet[u] = true
	}

	for i := uint64(1); i <= 128; i++ {
		host, err := client.getVolumeHost(ctx, i)
		if err == nil {
			volumes[i] = true
			assert.Equal(t, client.volumeIpMap[i], host)
			assert.Equal(t, hostSet[host], true)
			assert.Equal(t, err, nil)
		}
	}

	for i := uint64(1); i <= 128; i++ {
		host, err := client.getVolumeHost(ctx, i)
		if volumes[i] {
			assert.Equal(t, client.volumeIpMap[i], host)
			assert.Equal(t, hostSet[host], true)
			assert.Equal(t, err, nil)
		} else {
			assert.Equal(t, err, errors.ErrorCodeResponse[errors.ErrSeaweedFSMaster])
			assert.Equal(t, client.volumeIpMap[i], "")
		}
	}
}
