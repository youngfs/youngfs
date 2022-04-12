package ec

import (
	"github.com/go-playground/assert/v2"
	"icesos/set"
	"icesos/util"
	"testing"
)

func TestSetRules_IsLegal(t *testing.T) {
	//host is rand
	hostSet := make(map[string]bool)
	hosts := make([]string, 16)
	for i := range hosts {
		host := util.RandString(16)
		for hostSet[host] {
			host = util.RandString(16)
		}
		hosts[i] = host
	}

	setRules := &SetRules{
		Set:             set.Set(util.RandString(16)),
		Hosts:           hosts,
		DataShards:      uint64(len(hosts)) - 1,
		ParityShards:    1,
		MAXBlockSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), true)

	setRules = &SetRules{
		Set:             set.Set([]byte{0xff, 0xfe, 0xfd}),
		Hosts:           hosts,
		DataShards:      uint64(len(hosts)) - 1,
		ParityShards:    1,
		MAXBlockSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             set.Set(util.RandString(16)),
		Hosts:           hosts,
		DataShards:      uint64(len(hosts)),
		ParityShards:    1,
		MAXBlockSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             set.Set(util.RandString(16)),
		Hosts:           append(hosts, hosts...),
		DataShards:      uint64(len(hosts))*2 - 1,
		ParityShards:    1,
		MAXBlockSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             set.Set(util.RandString(16)),
		Hosts:           hosts,
		DataShards:      uint64(len(hosts)) - 1,
		ParityShards:    0,
		MAXBlockSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             set.Set(util.RandString(16)),
		Hosts:           hosts,
		DataShards:      uint64(len(hosts)),
		ParityShards:    0,
		MAXBlockSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             set.Set(util.RandString(16)),
		Hosts:           hosts,
		DataShards:      uint64(len(hosts)),
		ParityShards:    0,
		MAXBlockSize:    1024 * 1024 * 1024,
		ECMode:          false,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), true)

	setRules = &SetRules{
		Set:             set.Set(util.RandString(16)),
		Hosts:           hosts[0:1],
		DataShards:      uint64(len(hosts[0:1])),
		ParityShards:    0,
		MAXBlockSize:    1024 * 1024 * 1024,
		ECMode:          false,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             set.Set(util.RandString(16)),
		Hosts:           hosts[0:1],
		DataShards:      uint64(len(hosts[0:1])),
		ParityShards:    0,
		MAXBlockSize:    1024 * 1024 * 1024,
		ECMode:          false,
		ReplicationMode: false,
	}
	assert.Equal(t, setRules.IsLegal(), true)

	setRules = &SetRules{
		Set:             set.Set(util.RandString(16)),
		Hosts:           hosts[0:1],
		DataShards:      0,
		ParityShards:    0,
		MAXBlockSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             set.Set(util.RandString(16)),
		Hosts:           hosts[0:1],
		DataShards:      0,
		ParityShards:    0,
		MAXBlockSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             set.Set(util.RandString(16)),
		Hosts:           hosts,
		DataShards:      8,
		ParityShards:    0,
		MAXBlockSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             set.Set(util.RandString(16)),
		Hosts:           hosts,
		DataShards:      0,
		ParityShards:    8,
		MAXBlockSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             set.Set(util.RandString(16)),
		Hosts:           hosts,
		DataShards:      0,
		ParityShards:    0,
		MAXBlockSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             set.Set(util.RandString(16)),
		Hosts:           []string{},
		DataShards:      0,
		ParityShards:    0,
		MAXBlockSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: false,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             set.Set(util.RandString(16)),
		Hosts:           []string{},
		DataShards:      0,
		ParityShards:    0,
		MAXBlockSize:    1024 * 1024 * 1024,
		ECMode:          false,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), true)

	setRules = &SetRules{
		Set:             set.Set(util.RandString(16)),
		Hosts:           []string{},
		DataShards:      0,
		ParityShards:    0,
		MAXBlockSize:    1024 * 1024 * 1024,
		ECMode:          false,
		ReplicationMode: false,
	}
	assert.Equal(t, setRules.IsLegal(), true)

	setRules = &SetRules{
		Set:             set.Set(util.RandString(16)),
		Hosts:           nil,
		DataShards:      0,
		ParityShards:    0,
		MAXBlockSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: false,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             set.Set(util.RandString(16)),
		Hosts:           nil,
		DataShards:      0,
		ParityShards:    0,
		MAXBlockSize:    1024 * 1024 * 1024,
		ECMode:          false,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), true)

	setRules = &SetRules{
		Set:             set.Set(util.RandString(16)),
		Hosts:           nil,
		DataShards:      0,
		ParityShards:    0,
		MAXBlockSize:    1024 * 1024 * 1024,
		ECMode:          false,
		ReplicationMode: false,
	}
	assert.Equal(t, setRules.IsLegal(), true)

	hosts = make([]string, 0)

}
