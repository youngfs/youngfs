package entry

import (
	"encoding/hex"
	"os"
	"youngfs/fs/fullpath"
)

type ListEntry struct {
	FullPath string      // file full fullpath
	Set      string      // own bucket
	Mtime    string      // time of last modification
	Ctime    string      // time of creation
	Mode     os.FileMode // file mode
	Mime     string      // MIME type
	Md5      string      // MD5
	FileSize uint64      // file size
	Chunks   []ListChunk
}

type ListChunk struct {
	Offset        uint64     // offset
	Size          uint64     // size
	Md5           string     // MD5
	IsReplication bool       // is replication
	Frags         []ListFrag // frags
}

type ListFrag struct {
	Size        uint64 // size
	Id          int64  // id
	Md5         string // MD5
	IsDataShard bool   // is data shard
	Fid         string // fid
}

func (ent *Entry) ToListEntry() *ListEntry {
	if ent == nil {
		return nil
	}

	chunks := make([]ListChunk, len(ent.Chunks))
	for i, u := range ent.Chunks {
		chunks[i] = *u.ToListChunk()
	}

	return &ListEntry{
		FullPath: string(ent.FullPath),
		Set:      string(ent.Bucket),
		Mtime:    ent.Mtime.Format(timeFormat),
		Ctime:    ent.Ctime.Format(timeFormat),
		Mode:     ent.Mode,
		Mime:     ent.Mime,
		Md5:      hex.EncodeToString(ent.Md5),
		FileSize: ent.FileSize,
		Chunks:   chunks,
	}
}

func (c *Chunk) ToListChunk() *ListChunk {
	if c == nil {
		return nil
	}

	frags := make([]ListFrag, len(c.Frags))
	for i, u := range c.Frags {
		frags[i] = *u.ToListFrag()
	}

	return &ListChunk{
		Offset:        c.Offset,
		Size:          c.Size,
		Md5:           hex.EncodeToString(c.Md5),
		IsReplication: c.IsReplication,
		Frags:         frags,
	}
}

func (f *Frag) ToListFrag() *ListFrag {
	if f == nil {
		return nil
	}

	return &ListFrag{
		Size:        f.Size,
		Id:          f.Id,
		Md5:         hex.EncodeToString(f.Md5),
		IsDataShard: f.IsDataShard,
		Fid:         f.Fid,
	}
}

func ToListEntries(ents []Entry) []ListEntry {
	ret := make([]ListEntry, len(ents))

	for i, u := range ents {
		ret[i] = *u.ToListEntry()
	}

	return ret
}

func (ent *ListEntry) IsDirectory() bool {
	return ent.Mode.IsDir()
}

func (ent *ListEntry) IsFile() bool {
	return ent.Mode.IsRegular()
}

func (ent *ListEntry) Name() string {
	return fullpath.FullPath(ent.FullPath).Name()
}
