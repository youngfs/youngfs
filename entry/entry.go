package entry

import (
	"object-storage-server/util"
	"os"
	"time"
)

type Attribute struct {
	time     time.Time   // time of creation
	Mode     os.FileMode // file mode
	Mime     string      // MIME type
	TtlSec   uint64      // ttl in seconds
	Md5      []byte      // MD5
	FileSize uint64      // file size
}

func (attr Attribute) IsDirectory() bool {
	return attr.Mode&os.ModeDir > 0
}

type Entry struct {
	util.FullPath // file full path
	Attribute     // attribute
}

func (entry *Entry) Size() uint64 {
	return entry.FileSize
}

func (entry *Entry) TimeStamp() time.Time {
	return entry.time
}

func (entry *Entry) ShallowClone() *Entry {
	if entry == nil {
		return nil
	}
	newEntry := &Entry{}
	newEntry.FullPath = entry.FullPath
	newEntry.Attribute = entry.Attribute
	return newEntry
}
