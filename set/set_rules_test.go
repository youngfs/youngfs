package set

import (
	"github.com/go-playground/assert/v2"
	jsoniter "github.com/json-iterator/go"
	"icesfs/util"
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
		Set:             Set(util.RandString(16)),
		Hosts:           hosts,
		DataShards:      uint64(len(hosts)) - 1,
		ParityShards:    1,
		MAXShardSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), true)

	setRules = &SetRules{
		Set:             Set([]byte{0xff, 0xfe, 0xfd}),
		Hosts:           hosts,
		DataShards:      uint64(len(hosts)) - 1,
		ParityShards:    1,
		MAXShardSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             Set(util.RandString(16)),
		Hosts:           hosts,
		DataShards:      uint64(len(hosts)),
		ParityShards:    1,
		MAXShardSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             Set(util.RandString(16)),
		Hosts:           append(hosts, hosts...),
		DataShards:      uint64(len(hosts))*2 - 1,
		ParityShards:    1,
		MAXShardSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             Set(util.RandString(16)),
		Hosts:           hosts,
		DataShards:      uint64(len(hosts)) - 1,
		ParityShards:    0,
		MAXShardSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             Set(util.RandString(16)),
		Hosts:           hosts,
		DataShards:      uint64(len(hosts)),
		ParityShards:    0,
		MAXShardSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             Set(util.RandString(16)),
		Hosts:           hosts,
		DataShards:      uint64(len(hosts)),
		ParityShards:    0,
		MAXShardSize:    1024 * 1024 * 1024,
		ECMode:          false,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), true)

	setRules = &SetRules{
		Set:             Set(util.RandString(16)),
		Hosts:           hosts[0:1],
		DataShards:      uint64(len(hosts[0:1])),
		ParityShards:    0,
		MAXShardSize:    1024 * 1024 * 1024,
		ECMode:          false,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             Set(util.RandString(16)),
		Hosts:           hosts[0:1],
		DataShards:      uint64(len(hosts[0:1])),
		ParityShards:    0,
		MAXShardSize:    1024 * 1024 * 1024,
		ECMode:          false,
		ReplicationMode: false,
	}
	assert.Equal(t, setRules.IsLegal(), true)

	setRules = &SetRules{
		Set:             Set(util.RandString(16)),
		Hosts:           hosts[0:1],
		DataShards:      0,
		ParityShards:    0,
		MAXShardSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             Set(util.RandString(16)),
		Hosts:           hosts[0:1],
		DataShards:      0,
		ParityShards:    0,
		MAXShardSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             Set(util.RandString(16)),
		Hosts:           hosts,
		DataShards:      8,
		ParityShards:    0,
		MAXShardSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             Set(util.RandString(16)),
		Hosts:           hosts,
		DataShards:      0,
		ParityShards:    8,
		MAXShardSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             Set(util.RandString(16)),
		Hosts:           hosts,
		DataShards:      0,
		ParityShards:    0,
		MAXShardSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             Set(util.RandString(16)),
		Hosts:           []string{},
		DataShards:      0,
		ParityShards:    0,
		MAXShardSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: false,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             Set(util.RandString(16)),
		Hosts:           []string{},
		DataShards:      0,
		ParityShards:    0,
		MAXShardSize:    1024 * 1024 * 1024,
		ECMode:          false,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             Set(util.RandString(16)),
		Hosts:           []string{},
		DataShards:      0,
		ParityShards:    0,
		MAXShardSize:    1024 * 1024 * 1024,
		ECMode:          false,
		ReplicationMode: false,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             Set(util.RandString(16)),
		Hosts:           nil,
		DataShards:      0,
		ParityShards:    0,
		MAXShardSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: false,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             Set(util.RandString(16)),
		Hosts:           nil,
		DataShards:      0,
		ParityShards:    0,
		MAXShardSize:    1024 * 1024 * 1024,
		ECMode:          false,
		ReplicationMode: true,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	setRules = &SetRules{
		Set:             Set(util.RandString(16)),
		Hosts:           nil,
		DataShards:      0,
		ParityShards:    0,
		MAXShardSize:    1024 * 1024 * 1024,
		ECMode:          false,
		ReplicationMode: false,
	}
	assert.Equal(t, setRules.IsLegal(), false)

	hosts = make([]string, 0)

}

func TestSetRules_Json(t *testing.T) {
	//host is rand
	hostSet := make(map[string]bool)
	hosts := make([]string, 4)
	for i := range hosts {
		host := util.RandString(16)
		for hostSet[host] {
			host = util.RandString(16)
		}
		hosts[i] = host
	}

	setRules := &SetRules{
		Set:             Set(util.RandString(16)),
		Hosts:           hosts,
		DataShards:      uint64(len(hosts)) - 1,
		ParityShards:    1,
		MAXShardSize:    1024 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}

	val, err := jsoniter.Marshal(setRules)
	assert.Equal(t, err, nil)

	setRules2 := &SetRules{}
	err = jsoniter.Unmarshal(val, setRules2)
	assert.Equal(t, err, nil)
	assert.Equal(t, setRules2, setRules)
}
