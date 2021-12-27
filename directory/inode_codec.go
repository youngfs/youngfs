package directory

import (
	"github.com/golang/protobuf/proto"
	"icesos/directory/directory_pb"
	"icesos/full_path"
	"icesos/iam"
	"os"
	"time"
)

func (inode *Inode) toPb() *directory_pb.Inode {
	if inode == nil {
		return nil
	}
	return &directory_pb.Inode{
		FullPath: string(inode.FullPath),
		Set:      string(inode.Set),
		Time:     inode.Time.Unix(),
		Mode:     uint32(inode.Mode),
	}
}

func inodePbToInstance(pb *directory_pb.Inode) *Inode {
	if pb == nil {
		return nil
	}
	return &Inode{
		FullPath: full_path.FullPath(pb.FullPath),
		Set:      iam.Set(pb.Set),
		Time:     time.Unix(pb.Time, 0),
		Mode:     os.FileMode(pb.Mode),
	}
}

func (inode *Inode) encodeProto() ([]byte, error) {
	message := inode.toPb()
	return proto.Marshal(message)
}

func decodeInodeProto(b []byte) (*Inode, error) {
	message := &directory_pb.Inode{}
	if err := proto.Unmarshal(b, message); err != nil {
		return nil, err
	}
	return inodePbToInstance(message), nil
}
