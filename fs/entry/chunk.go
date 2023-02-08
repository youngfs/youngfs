package entry

import "sort"

type Chunk struct {
	Offset        uint64 // offset
	Size          uint64 // size
	Md5           []byte // MD5
	IsReplication bool   // is replication
	Frags                // frags
}

type Chunks []*Chunk

func (c Chunks) Len() int           { return len(c) }
func (c Chunks) Less(i, j int) bool { return c[i].Offset < c[j].Offset }
func (c Chunks) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

func (c Chunks) Verify() bool {
	sort.Sort(c)
	offset := uint64(0)
	if len(c) == 0 {
		return false
	}
	for _, chunk := range c {
		if chunk == nil {
			return false
		}
		if chunk.Offset != offset {
			return false
		}
		offset += chunk.Size
		if len(chunk.Frags) == 0 {
			return false
		}
		sort.Sort(chunk.Frags)
		fragSize := uint64(0)
		for i, frag := range chunk.Frags {
			if frag == nil {
				return false
			}
			if frag.Id != int64(i)+1 {
				return false
			}
			if chunk.IsReplication && chunk.Size != frag.Size {
				return false
			}
			if frag.IsDataShard {
				fragSize += frag.Size
			}
		}
		if (!chunk.IsReplication && fragSize != chunk.Size) || fragSize == 0 {
			return false
		}
	}
	return true
}
