package entry

type Frag struct {
	Size          uint64 // size
	Id            int64  // id
	Md5           []byte // MD5
	IsReplication bool   // is replication
	IsDataShard   bool   // is data shard
	Fid           string // fid
}

type Frags []Frag

func (f Frags) Len() int           { return len(f) }
func (f Frags) Less(i, j int) bool { return f[i].Id < f[j].Id }
func (f Frags) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
