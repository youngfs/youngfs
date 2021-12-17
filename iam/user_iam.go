package iam

import (
	"crypto/md5"
	"icesos/iam/iam_pb"
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

	skMd5 := make([]byte, 16)
	for i, b := range userIam.SecretKey {
		skMd5[i] = b
	}

	return &iam_pb.UserIAM{
		User:      string(userIam.User),
		SecretKey: skMd5,
	}
}

func userIAMPbToInstance(pb *iam_pb.UserIAM) *UserIAM {
	if pb == nil || len(pb.SecretKey) != md5.Size {
		return nil
	}

	var skMd5 [md5.Size]byte
	for i, b := range pb.SecretKey {
		skMd5[i] = b
	}

	return &UserIAM{
		User:      User(pb.User),
		SecretKey: skMd5,
	}
}
