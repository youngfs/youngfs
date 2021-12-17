package iam

import (
	"github.com/golang/protobuf/proto"
	"icesos/iam/iam_pb"
)

func (userIam *UserIAM) EncodeProto() ([]byte, error) {
	message := userIam.toPb()
	return proto.Marshal(message)
}

func DecodeUserIAMProto(b []byte) (*UserIAM, error) {
	message := &iam_pb.UserIAM{}
	if err := proto.Unmarshal(b, message); err != nil {
		return nil, err
	}
	return userIAMPbToInstance(message), nil
}
