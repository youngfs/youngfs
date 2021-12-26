package entry

import (
	"crypto/md5"
	"icesos/full_path"
	"icesos/iam"
	"icesos/kv"
	"os"
	"time"
)

type Entry struct {
	full_path.FullPath                // file full full_path
	iam.Set                           // own set_iam
	Time               time.Time      // time of creation
	Mode               os.FileMode    // file mode
	Mime               string         // MIME type
	Md5                [md5.Size]byte // MD5
	FileSize           uint64         // file size
	VolumeId           uint64         // volume id
	Fid                string         // fid
}

func (entry *Entry) key() string {
	return string(entry.Set) + "_" + string(entry.FullPath) + entryKv
}

func entryKey(fp full_path.FullPath, set iam.Set) string {
	return string(set) + "_" + string(fp) + entryKv
}

func InsertEntry(entry *Entry) error {
	return nil
}

func GetEntry(fp full_path.FullPath, set iam.Set) (*Entry, error) {
	key := entryKey(fp, set)

	b, err := kv.Client.KvGet(key)
	if err != nil {
		return nil, err
	}

	return decodeEntryProto(b)
}
