package set_iam

import (
	"github.com/golang/protobuf/proto"
	"icesos/iam/iam_pb"
)

func (setIam *SetIAM) EncodeProto() ([]byte, error) {
	message := setIam.toPb()
	return proto.Marshal(message)
}

func DecodeSetIAMProto(b []byte) (*SetIAM, error) {
	message := &iam_pb.SetIAM{}
	if err := proto.Unmarshal(b, message); err != nil {
		return nil, err
	}
	return setIAMPbToInstance(message), nil
}
