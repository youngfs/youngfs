package directory

import (
	"github.com/golang/protobuf/proto"
	"icesos/directory/directory_pb"
)

func (inode *Inode) EncodeProto() ([]byte, error) {
	message := inode.toPb()
	return proto.Marshal(message)
}

func DecodeInodeProto(b []byte) (*Inode, error) {
	message := &directory_pb.Inode{}
	if err := proto.Unmarshal(b, message); err != nil {
		return nil, err
	}
	return inodePbToInstance(message), nil
}
