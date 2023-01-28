package entry

import (
	"os"
	"time"
	"youngfs/fs/full_path"
	"youngfs/fs/set"
)

type Entry struct {
	full_path.FullPath             // file full full_path
	set.Set                        // own set
	Mtime              time.Time   // time of last modification
	Ctime              time.Time   // time of creation
	Mode               os.FileMode // file mode
	Mime               string      // MIME type
	Md5                []byte      // MD5
	FileSize           uint64      // file size
	Fid                string      // fid
	ECid               string      // erasure code id
}

func (ent *Entry) Key() string {
	return string(ent.Set) + string(ent.FullPath) + entryKv
}

func EntryKey(set set.Set, fp full_path.FullPath) string {
	return string(set) + string(fp) + entryKv
}

func (ent *Entry) IsDirectory() bool {
	return ent.Mode.IsDir()
}

func (ent *Entry) IsFile() bool {
	return ent.Mode.IsRegular()
}
