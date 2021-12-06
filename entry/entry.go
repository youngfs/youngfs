package entry

import (
	"icesos/entry/entry_pb"
	"icesos/util"
	"os"
	"time"
)

type Attribute struct {
	Time     time.Time   // time of creation
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
	*Attribute    // attribute
}

func (entry *Entry) Size() uint64 {
	return entry.FileSize
}

func (entry *Entry) TimeStamp() time.Time {
	return entry.Time
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

func (entry *Entry) ToPb() *entry_pb.Entry {
	if entry == nil {
		return nil
	}
	return &entry_pb.Entry{
		FullPath:    string(entry.FullPath),
		IsDirectory: entry.IsDirectory(),
		Attribute:   entry.Attribute.ToPb(),
	}
}

func (attr *Attribute) ToPb() *entry_pb.Attribute {
	if attr == nil {
		return nil
	}
	return &entry_pb.Attribute{
		Time:     attr.Time.Unix(),
		Mode:     uint32(attr.Mode),
		Mine:     attr.Mime,
		TtlSec:   attr.TtlSec,
		Md5:      attr.Md5,
		FileSize: attr.FileSize,
	}
}

func EntryPbToInstance(pb *entry_pb.Entry) *Entry {
	if pb == nil {
		return nil
	}
	return &Entry{
		FullPath:  util.FullPath(pb.FullPath),
		Attribute: AttributePbTonstance(pb.Attribute),
	}
}

func AttributePbTonstance(pb *entry_pb.Attribute) *Attribute {
	if pb == nil {
		return nil
	}
	return &Attribute{
		Time:     time.Unix(pb.Time, 0),
		Mode:     os.FileMode(pb.Mode),
		Mime:     pb.Mine,
		TtlSec:   pb.TtlSec,
		Md5:      pb.Md5,
		FileSize: pb.FileSize,
	}
}
