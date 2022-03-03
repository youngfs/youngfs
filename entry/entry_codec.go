package entry

import (
	"crypto/md5"
	"github.com/golang/protobuf/proto"
	"icesos/entry/entry_pb"
	"icesos/errors"
	"icesos/full_path"
	"icesos/set"
	"icesos/util"
	"os"
	"time"
)

func (ent *Entry) ToPb() *entry_pb.Entry {
	if ent == nil {
		return nil
	}

	return &entry_pb.Entry{
		FullPath: string(ent.FullPath),
		Set:      string(ent.Set),
		Ctime:    ent.Ctime.Unix(),
		Mode:     uint32(ent.Mode),
		Mine:     ent.Mime,
		Md5:      util.Md5ToBytes(ent.Md5),
		FileSize: ent.FileSize,
		Fid:      ent.Fid,
	}
}

func EntryPbToInstance(pb *entry_pb.Entry) *Entry {
	if pb == nil || len(pb.Md5) != md5.Size {
		return nil
	}

	return &Entry{
		FullPath: full_path.FullPath(pb.FullPath),
		Set:      set.Set(pb.Set),
		Ctime:    time.Unix(pb.Ctime, 0),
		Mode:     os.FileMode(pb.Mode),
		Mime:     pb.Mine,
		Md5:      util.BytesToMd5(pb.Md5),
		FileSize: pb.FileSize,
		Fid:      pb.Fid,
	}
}

func (ent *Entry) EncodeProto() ([]byte, error) {
	message := ent.ToPb()
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
	return EntryPbToInstance(message), nil
}
