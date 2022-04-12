package ec

import (
	"icesos/full_path"
	"icesos/set"
)

type Frag struct {
	full_path.FullPath        // full path
	set.Set                   // set
	Fid                string // fid
	FileSize           uint64 // file size
}

//Frags must not be []Frag{} (codec will become nill)
type Shard struct {
	Host  string // host
	Frags []Frag // frags
}

//DataBlocks must not be []Shard{} (codec will become nill)
type Suite struct {
	ECid   string  // erasure code id
	Fid    string  // fid
	Shards []Shard // data blocks
}
