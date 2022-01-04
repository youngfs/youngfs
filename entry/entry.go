package entry

import (
	"crypto/md5"
	"icesos/full_path"
	"icesos/iam"
	"icesos/kv"
	"icesos/storage_engine"
	"os"
	"time"
)

type Entry struct {
	full_path.FullPath                // file full full_path
	iam.Set                           // own set_iam
	Mtime              time.Time      // time of last modification
	Ctime              time.Time      // time of creation
	Mode               os.FileMode    // file mode
	Mime               string         // MIME type
	Md5                [md5.Size]byte // MD5
	FileSize           uint64         // file size
	VolumeId           uint64         // volume id
	Fid                string         // fid
}

func (entry *Entry) key() string {
	return string(entry.Set) + string(entry.FullPath) + entryKv
}

func entryKey(set iam.Set, fp full_path.FullPath) string {
	return string(set) + string(fp) + entryKv
}

func (entry *Entry) IsDirectory() bool {
	return entry.Mode.IsDir()
}

func (entry *Entry) IsFile() bool {
	return entry.Mode.IsRegular()
}

// after insert entry, insert inode
func InsertEntry(entry *Entry) error {
	b, err := entry.encodeProto()
	if err != nil {
		return err
	}

	return kv.Client.KvPut(entry.key(), b)
}

func GetEntry(set iam.Set, fp full_path.FullPath) (*Entry, error) {
	key := entryKey(set, fp)

	b, err := kv.Client.KvGet(key)
	if err != nil {
		return nil, err
	}

	return decodeEntryProto(b)
}

// after delete entry, delete inode
func DeleteEntry(set iam.Set, fp full_path.FullPath) error {
	key := entryKey(set, fp)

	entry, err := GetEntry(set, fp)
	if err != nil {
		return err
	}
	_, err = kv.Client.KvDelete(key)
	if err != nil {
		return err
	}

	if entry.IsFile() {
		err := storage_engine.DeleteObject(entry.VolumeId, entry.Fid, entry.FileSize)
		if err != nil {
			return err
		}
	}

	return nil
}
