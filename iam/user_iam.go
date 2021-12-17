package iam

import (
	"crypto/md5"
	"icesos/iam/iam_pb"
	"icesos/util"
)

type UserIAM struct {
	User                     // user_iam name
	SecretKey [md5.Size]byte // secret key md5
}

func (userIam *UserIAM) Key() string {
	return string(userIam.User) + userIAMKey
}

func (userIam *UserIAM) toPb() *iam_pb.UserIAM {
	if userIam == nil {
		return nil
	}

	return &iam_pb.UserIAM{
		User:      string(userIam.User),
		SecretKey: util.Md5ToBytes(userIam.SecretKey),
	}
}

func userIAMPbToInstance(pb *iam_pb.UserIAM) *UserIAM {
	if pb == nil || len(pb.SecretKey) != md5.Size {
		return nil
	}

	return &UserIAM{
		User:      User(pb.User),
		SecretKey: util.BytesToMd5(pb.SecretKey),
	}
}
