package rules

import "youngfs/fs/set"

type Rules struct {
	Set             set.Set  `form:"Set" json:"Set" uri:"Set" xml:"Set" yaml:"Set"`                                                             // set
	Hosts           []string `form:"Hosts" json:"Hosts" uri:"Hosts" xml:"Hosts" yaml:"Hosts"`                                                   // hosts len(Hosts) >= DataShards + ParityShards
	DataShards      uint64   `form:"DataShards" json:"DataShards" uri:"DataShards" xml:"DataShards" yaml:"DataShards"`                          // business data block number
	ParityShards    uint64   `form:"ParityShards" json:"ParityShards" uri:"ParityShards" xml:"ParityShards" yaml:"ParityShards"`                // redundant data block number
	MaxShardSize    uint64   `form:"MaxShardSize" json:"MaxShardSize" uri:"MaxShardSize" xml:"MaxShardSize" yaml:"MaxShardSize"`                // max shard size
	ECMode          bool     `form:"ECMode" json:"ECMode" uri:"ECMode" xml:"ECMode" yaml:"ECMode"`                                              // is EC mode on
	ReplicationMode bool     `form:"ReplicationMode" json:"ReplicationMode" uri:"ReplicationMode" xml:"ReplicationMode" yaml:"ReplicationMode"` // is Replication mode on
}

// still need to check hosts is legal
func (rules *Rules) IsLegal() bool {
	if !rules.Set.IsLegal() {
		return false
	}

	if rules.Hosts == nil || len(rules.Hosts) == 0 {
		return false
	}

	if rules.ECMode {
		if rules.DataShards+rules.ParityShards > 32 {
			return false
		}

		if uint64(len(rules.Hosts)) < rules.DataShards+rules.ParityShards {
			return false
		}

		if rules.ParityShards < 1 || rules.DataShards < 1 {
			return false
		}
	}

	if rules.ReplicationMode {
		if len(rules.Hosts) < 2 {
			return false
		}
	}

	hostSet := make(map[string]bool)
	for _, u := range rules.Hosts {
		if hostSet[u] {
			return false
		}
		hostSet[u] = true
	}

	return true
}
