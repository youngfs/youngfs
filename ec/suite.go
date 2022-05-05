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
	Md5   []byte // MD5
}

//DataBlocks must not be []Shard{} (codec will become nill)
type Suite struct {
	ECid               string  // erasure code id
	full_path.FullPath         // full path
	set.Set                    // set
	OrigFid            string  // original fid
	FileSize           uint64  // file size
	BakHost            string  // backup host
	BakFid             string  // backup fid
	Next               string  // next ECid, if it's end,next = ""
	Shards             []Shard // data blocks
}
