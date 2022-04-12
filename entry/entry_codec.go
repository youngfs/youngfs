package entry

import (
	"github.com/golang/protobuf/proto"
	"icesos/entry/entry_pb"
	"icesos/errors"
	"icesos/full_path"
	"icesos/set"
	"os"
	"time"
)

func (ent *Entry) toPb() *entry_pb.Entry {
	if ent == nil {
		return nil
	}

	return &entry_pb.Entry{
		FullPath: string(ent.FullPath),
		Set:      string(ent.Set),
		Mtime:    ent.Mtime.Unix(),
		Ctime:    ent.Ctime.Unix(),
		Mode:     uint32(ent.Mode),
		Mine:     ent.Mime,
		Md5:      ent.Md5,
		FileSize: ent.FileSize,
		Fid:      ent.Fid,
		ECid:     ent.ECid,
	}
}

func entryPbToInstance(pb *entry_pb.Entry) *Entry {
	if pb == nil {
		return nil
	}

	return &Entry{
		FullPath: full_path.FullPath(pb.FullPath),
		Set:      set.Set(pb.Set),
		Mtime:    time.Unix(pb.Mtime, 0),
		Ctime:    time.Unix(pb.Ctime, 0),
		Mode:     os.FileMode(pb.Mode),
		Mime:     pb.Mine,
		Md5:      pb.Md5,
		FileSize: pb.FileSize,
		Fid:      pb.Fid,
		ECid:     pb.ECid,
	}
}

func (ent *Entry) EncodeProto() ([]byte, error) {
	message := ent.toPb()
	b, err := proto.Marshal(message)
	if err != nil {
		err = errors.ErrorCodeResponse[errors.ErrProto]
	}
	return b, err
}

func DecodeEntryProto(b []byte) (*Entry, error) {
	message := &entry_pb.Entry{}
	if err := proto.Unmarshal(b, message); err != nil {
		return nil, errors.ErrorCodeResponse[errors.ErrProto]
	}
	return entryPbToInstance(message), nil
}
