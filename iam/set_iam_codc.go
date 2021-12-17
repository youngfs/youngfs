package iam

import (
	"github.com/golang/protobuf/proto"
	"icesos/iam/iam_pb"
)

func (setIam *setIAM) toPb() *iam_pb.SetIAM {
	if setIam == nil {
		return nil
	}
	return &iam_pb.SetIAM{
		User: string(setIam.User),
		Set:  string(setIam.Set),
	}
}

func setIAMPbToInstance(pb *iam_pb.SetIAM) *setIAM {
	if pb == nil {
		return nil
	}
	return &setIAM{
		User: User(pb.User),
		Set:  Set(pb.Set),
	}
}

func (setIam *setIAM) encodeProto() ([]byte, error) {
	message := setIam.toPb()
	return proto.Marshal(message)
}

func decodeSetIAMProto(b []byte) (*setIAM, error) {
	message := &iam_pb.SetIAM{}
	if err := proto.Unmarshal(b, message); err != nil {
		return nil, err
	}
	return setIAMPbToInstance(message), nil
}
