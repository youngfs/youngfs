package object

import (
	"crypto/md5"
	"github.com/golang/protobuf/proto"
	"icesos/full_path"
	"icesos/iam"
	"icesos/object/object_pb"
	"icesos/util"
	"os"
	"time"
)

func (ob *Object) toPb() *object_pb.Object {
	if ob == nil {
		return nil
	}

	return &object_pb.Object{
		FullPath: string(ob.FullPath),
		Set:      string(ob.Set),
		Time:     ob.Time.Unix(),
		Mode:     uint32(ob.Mode),
		Mine:     ob.Mime,
		Md5:      util.Md5ToBytes(ob.Md5),
		FileSize: ob.FileSize,
		VolumeId: ob.VolumeId,
		Fid:      ob.Fid,
	}
}

func objectPbToInstance(pb *object_pb.Object) *Object {
	if pb == nil || len(pb.Md5) != md5.Size {
		return nil
	}

	return &Object{
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

func (ob *Object) encodeProto() ([]byte, error) {
	message := ob.toPb()
	return proto.Marshal(message)
}

func decodeObjectProto(b []byte) (*Object, error) {
	message := &object_pb.Object{}
	if err := proto.Unmarshal(b, message); err != nil {
		return nil, err
	}
	return objectPbToInstance(message), nil
}
