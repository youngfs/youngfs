package object

import (
	"github.com/golang/protobuf/proto"
	"icesos/object/object_pb"
)

func (ob *Object) EncodeProto() ([]byte, error) {
	message := ob.toPb()
	return proto.Marshal(message)
}

func DecodeObjectProto(b []byte) (*Object, error) {
	message := &object_pb.Object{}
	if err := proto.Unmarshal(b, message); err != nil {
		return nil, err
	}
	return objectPbToInstance(message), nil
}
