package entry

import (
	"crypto/md5"
	"github.com/golang/protobuf/proto"
	"icesos/entry/entry_pb"
	"icesos/full_path"
	"icesos/iam"
	"icesos/util"
	"os"
	"time"
)

func (entry *Entry) toPb() *entry_pb.Entry {
	if entry == nil {
		return nil
	}

	return &entry_pb.Entry{
		FullPath: string(entry.FullPath),
		Set:      string(entry.Set),
		Time:     entry.Time.Unix(),
		Mode:     uint32(entry.Mode),
		Mine:     entry.Mime,
		Md5:      util.Md5ToBytes(entry.Md5),
		FileSize: entry.FileSize,
		VolumeId: entry.VolumeId,
		Fid:      entry.Fid,
	}
}

func entryPbToInstance(pb *entry_pb.Entry) *Entry {
	if pb == nil || len(pb.Md5) != md5.Size {
		return nil
	}

	return &Entry{
		FullPath: full_path.FullPath(pb.FullPath),
		Set:      iam.Set(pb.Set),
		Time:     time.Unix(pb.Time, 0),
		Mode:     os.FileMode(pb.Mode),
		Mime:     pb.Mine,
		Md5:      util.BytesToMd5(pb.Md5),
		FileSize: pb.FileSize,
		VolumeId: pb.VolumeId,
		Fid:      pb.Fid,
	}
}

func (entry *Entry) encodeProto() ([]byte, error) {
	message := entry.toPb()
	return proto.Marshal(message)
}

func decodeEntryProto(b []byte) (*Entry, error) {
	message := &entry_pb.Entry{}
	if err := proto.Unmarshal(b, message); err != nil {
		return nil, err
	}
	return entryPbToInstance(message), nil
}
