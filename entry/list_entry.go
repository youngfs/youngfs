package entry

import (
	"encoding/hex"
	"icesos/full_path"
	"os"
)

type ListEntry struct {
	FullPath string      // file full full_path
	Set      string      // own set
	Mtime    string      // time of last modification
	Ctime    string      // time of creation
	Mode     os.FileMode // file mode
	Mime     string      // MIME type
	Md5      string      // MD5
	FileSize uint64      // file size
	Fid      string      // fid
	ECid     string      // erasure code id
}

func (ent *Entry) ToListEntry() *ListEntry {
	return &ListEntry{
		FullPath: string(ent.FullPath),
		Set:      string(ent.Set),
		Mtime:    ent.Mtime.Format(timeFormat),
		Ctime:    ent.Ctime.Format(timeFormat),
		Mode:     ent.Mode,
		Mime:     ent.Mime,
		Md5:      hex.EncodeToString(ent.Md5),
		FileSize: ent.FileSize,
		Fid:      ent.Fid,
		ECid:     ent.ECid,
	}
}

func ToListEntris(ents []Entry) []ListEntry {
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
	return full_path.FullPath(ent.FullPath).Name()
}
