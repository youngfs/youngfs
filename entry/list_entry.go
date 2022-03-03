package entry

import (
	"icesos/full_path"
	"icesos/set"
	"icesos/util"
	"os"
)

type ListEntry struct {
	FullPath full_path.FullPath // file full full_path
	Set      set.Set            // own set
	Mode     os.FileMode        // file mode
	Md5      string             // MD5
	FileSize uint64             // file size
	Fid      string             // fid
}

func (ent *Entry) ToListEntry() *ListEntry {
	return &ListEntry{
		FullPath: ent.FullPath,
		Set:      ent.Set,
		Mode:     ent.Mode,
		Md5:      util.Md5ToStr(ent.Md5),
		FileSize: ent.FileSize,
		Fid:      ent.Fid,
	}
}

func ToListEntris(ents []Entry) []ListEntry {
	ret := make([]ListEntry, len(ents))

	for i, u := range ents {
		ret[i] = *u.ToListEntry()
	}

	return ret
}
