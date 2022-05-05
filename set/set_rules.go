package set

type SetRules struct {
	Set             `json:"Set"` // set
	Hosts           []string     `json:"Hosts"`           // hosts len(Hosts) >= DataShards + ParityShards
	DataShards      uint64       `json:"DataShards"`      // business data block number
	ParityShards    uint64       `json:"ParityShards"`    // redundant data block number
	MAXShardSize    uint64       `json:"MAXShardSize"`    // max shard size
	ECMode          bool         `json:"ECMode"`          // is EC mode on
	ReplicationMode bool         `json:"ReplicationMode"` // is Replication mode on
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
		if setRules.DataShards+setRules.ParityShards > 256 {
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
