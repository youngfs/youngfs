package vfs

import (
	"crypto/md5"
	"github.com/golang/protobuf/proto"
	"icesos/entry/entry_pb"
	"icesos/errors"
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
		Mtime:    entry.Mtime.Unix(),
		Ctime:    entry.Ctime.Unix(),
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
		Mtime:    time.Unix(pb.Mtime, 0),
		Ctime:    time.Unix(pb.Ctime, 0),
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
	b, err := proto.Marshal(message)
	if err != nil {
		err = errors.ErrorCodeResponse[errors.ErrProto]
	}
	return b, err
}

func decodeEntryProto(b []byte) (*Entry, error) {
	message := &entry_pb.Entry{}
	if err := proto.Unmarshal(b, message); err != nil {
		return nil, errors.ErrorCodeResponse[errors.ErrProto]
	}
	return entryPbToInstance(message), nil
}
