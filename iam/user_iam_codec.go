package iam

import (
	"crypto/md5"
	"github.com/golang/protobuf/proto"
	"icesos/errors"
	"icesos/iam/iam_pb"
	"icesos/util"
)

func (userIam *userIAM) toPb() *iam_pb.UserIAM {
	if userIam == nil {
		return nil
	}

	return &iam_pb.UserIAM{
		User:      string(userIam.User),
		SecretKey: util.Md5ToBytes(userIam.SecretKey),
	}
}

func userIAMPbToInstance(pb *iam_pb.UserIAM) *userIAM {
	if pb == nil || len(pb.SecretKey) != md5.Size {
		return nil
	}

	return &userIAM{
		User:      User(pb.User),
		SecretKey: util.BytesToMd5(pb.SecretKey),
	}
}

func (userIam *userIAM) encodeProto() ([]byte, error) {
	message := userIam.toPb()
	b, err := proto.Marshal(message)
	if err != nil {
		err = errors.ErrorCodeResponse[errors.ErrProto]
	}
	return b, err
}

func decodeUserIAMProto(b []byte) (*userIAM, error) {
	message := &iam_pb.UserIAM{}
	if err := proto.Unmarshal(b, message); err != nil {
		return nil, errors.ErrorCodeResponse[errors.ErrProto]
	}
	return userIAMPbToInstance(message), nil
}
