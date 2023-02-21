package seaweedfs

import (
	"context"
	"github.com/go-playground/assert/v2"
	"testing"
	"youngfs/errors"
	"youngfs/vars"
)

func TestSeaweedFS_GetVolumeHost(t *testing.T) {
	client := NewStorageEngine(vars.SeaweedFSMaster)

	hosts, err := client.GetHosts(context.Background())
	assert.Equal(t, err, nil)

	volumes := make(map[uint64]bool)
	hostSet := make(map[string]bool)
	for _, u := range hosts {
		hostSet[u] = true
	}

	for i := uint64(1); i <= 128; i++ {
		host, err := client.getVolumeHost(i)
		if err == nil {
			volumes[i] = true
			val, ok := client.volumeIpMap.Load(i)
			assert.Equal(t, ok, true)
			ip, ok := val.(string)
			assert.Equal(t, ok, true)
			assert.Equal(t, ip, host)
			assert.Equal(t, hostSet[host], true)
			assert.Equal(t, err, nil)
		}
	}

	for i := uint64(1); i <= 128; i++ {
		host, err := client.getVolumeHost(i)
		if volumes[i] {
			val, ok := client.volumeIpMap.Load(i)
			assert.Equal(t, ok, true)
			ip, ok := val.(string)
			assert.Equal(t, ok, true)
			assert.Equal(t, ip, host)
			assert.Equal(t, hostSet[host], true)
			assert.Equal(t, err, nil)
		} else {
			assert.Equal(t, errors.Is(err, errors.ErrSeaweedFSMaster), true)
			val, ok := client.volumeIpMap.Load(i)
			assert.Equal(t, ok, false)
			assert.Equal(t, val, nil)
		}
	}
}
