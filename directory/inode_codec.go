package directory

import (
	"github.com/golang/protobuf/proto"
	"icesos/directory/directory_pb"
	"icesos/errors"
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
		Mtime:    inode.Mtime.Unix(),
		Ctime:    inode.Ctime.Unix(),
		Mode:     uint32(inode.Mode),
		Mime:     inode.Mime,
		FileSize: inode.FileSize,
	}
}

func inodePbToInstance(pb *directory_pb.Inode) *Inode {
	if pb == nil {
		return nil
	}
	return &Inode{
		FullPath: full_path.FullPath(pb.FullPath),
		Set:      iam.Set(pb.Set),
		Mtime:    time.Unix(pb.Mtime, 0),
		Ctime:    time.Unix(pb.Ctime, 0),
		Mode:     os.FileMode(pb.Mode),
		Mime:     pb.Mime,
		FileSize: pb.FileSize,
	}
}

func (inode *Inode) encodeProto() ([]byte, error) {
	message := inode.toPb()
	b, err := proto.Marshal(message)
	if err != nil {
		err = errors.ErrorCodeResponse[errors.ErrProto]
	}
	return b, err
}

func decodeInodeProto(b []byte) (*Inode, error) {
	message := &directory_pb.Inode{}
	if err := proto.Unmarshal(b, message); err != nil {
		return nil, errors.ErrorCodeResponse[errors.ErrProto]
	}
	return inodePbToInstance(message), nil
}
