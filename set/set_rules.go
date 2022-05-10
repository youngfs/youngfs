package set

type SetRules struct {
	Set             Set      `form:"Set" json:"Set" uri:"Set" xml:"Set" yaml:"Set"`                                                             // set
	Hosts           []string `form:"Hosts" json:"Hosts" uri:"Hosts" xml:"Hosts" yaml:"Hosts"`                                                   // hosts len(Hosts) >= DataShards + ParityShards
	DataShards      uint64   `form:"DataShards" json:"DataShards" uri:"DataShards" xml:"DataShards" yaml:"DataShards"`                          // business data block number
	ParityShards    uint64   `form:"ParityShards" json:"ParityShards" uri:"ParityShards" xml:"ParityShards" yaml:"ParityShards"`                // redundant data block number
	MAXShardSize    uint64   `form:"MAXShardSize" json:"MAXShardSize" uri:"MAXShardSize" xml:"MAXShardSize" yaml:"MAXShardSize"`                // max shard size
	ECMode          bool     `form:"ECMode" json:"ECMode" uri:"ECMode" xml:"ECMode" yaml:"ECMode"`                                              // is EC mode on
	ReplicationMode bool     `form:"ReplicationMode" json:"ReplicationMode" uri:"ReplicationMode" xml:"ReplicationMode" yaml:"ReplicationMode"` // is Replication mode on
}

// still need to check hosts is legal
func (setRules *SetRules) IsLegal() bool {
	if !setRules.Set.IsLegal() {
		return false
	}

	if setRules.Hosts == nil || len(setRules.Hosts) == 0 {
		if setRules.ECMode {
			return false
		} else {
			return true
		}
	}

	if setRules.ECMode {
		if setRules.DataShards+setRules.ParityShards > 32 {
			return false
		}

		if uint64(len(setRules.Hosts)) < setRules.DataShards+setRules.ParityShards {
			return false
		}

		if setRules.ParityShards < 1 || setRules.DataShards < 1 {
			return false
		}
	}

	if setRules.ReplicationMode {
		if len(setRules.Hosts) < 2 {
			return false
		}
	}

	hostSet := make(map[string]bool)
	for _, u := range setRules.Hosts {
		if hostSet[u] {
			return false
		}
		hostSet[u] = true
	}

	return true
}
