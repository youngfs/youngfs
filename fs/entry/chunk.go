package entry

import "sort"

type Chunk struct {
	Offset uint64 // offset
	Size   uint64 // size
	Md5    []byte // MD5
	Frags         // frags
}

type Chunks []Chunk

func (c Chunks) Len() int           { return len(c) }
func (c Chunks) Less(i, j int) bool { return c[i].Offset < c[j].Offset }
func (c Chunks) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

func (c Chunks) Verify() bool {
	sort.Sort(c)
	offset := uint64(0)
	for _, chunk := range c {
		if chunk.Offset != offset {
			return false
		}
		offset += chunk.Size
	}
	return true
}
