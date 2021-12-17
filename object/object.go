package object

import (
	"icesos/full_path"
	"icesos/iam"
	"icesos/object/object_pb"
	"os"
	"time"
)

type Object struct {
	full_path.FullPath             // file full full_path
	set_iam.Set                    // own set_iam
	Time               time.Time   // time of creation
	Mode               os.FileMode // file mode
	Mime               string      // MIME type
	Md5                []byte      // MD5
	FileSize           uint64      // file size
	VolumeId           uint64      // volume id
	Fid                string      // fid
}

func (ob *Object) Key() string {
	return string(ob.Set) + "_" + string(ob.FullPath) + "_object"
}

func (ob *Object) TimeUnix() int64 {
	return ob.Time.Unix()
}

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
		Md5:      ob.Md5,
		FileSize: ob.FileSize,
		VolumeId: ob.VolumeId,
		Fid:      ob.Fid,
	}
}

func objectPbToInstance(pb *object_pb.Object) *Object {
	if pb == nil {
		return nil
	}
	return &Object{
		FullPath: full_path.FullPath(pb.FullPath),
		Set:      set_iam.Set(pb.Set),
		Time:     time.Unix(pb.Time, 0),
		Mode:     os.FileMode(pb.Mode),
		Mime:     pb.Mine,
		Md5:      pb.Md5,
		FileSize: pb.FileSize,
		VolumeId: pb.VolumeId,
		Fid:      pb.Fid,
	}
}
