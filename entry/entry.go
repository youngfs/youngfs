package entry

import (
	"crypto/md5"
	"icesos/full_path"
	"icesos/set"
	"os"
	"time"
)

type Entry struct {
	full_path.FullPath                // file full full_path
	set.Set                           // own set
	Ctime              time.Time      // time of creation
	Mode               os.FileMode    // file mode
	Mime               string         // MIME type
	Md5                [md5.Size]byte // MD5
	FileSize           uint64         // file size
	Fid                string         // fid
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
