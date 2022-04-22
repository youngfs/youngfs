package ec

import (
	"icesos/set"
)

type SetRules struct {
	set.Set                  // set
	Hosts           []string // hosts len(Hosts) >= DataShards + ParityShards
	DataShards      uint64   // business data block number
	ParityShards    uint64   // redundant data block number
	MAXBlockSize    uint64   // max block size
	ECMode          bool     // is EC mode on
	ReplicationMode bool     // is Replication mode on
}

func (setRules *SetRules) Key() string {
	return string(setRules.Set) + setRulesKey
}

func SetRulesKey(set set.Set) string {
	return string(set) + setRulesKey
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

func (ec *EC) InsertSetRules(setRules *SetRules) error {

	return nil
}
