package iam

import (
	"icesos/iam/iam_pb"
)

type UserIAM struct {
	User             // user_iam name
	SecretKey string // secret key
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
		SecretKey: userIam.SecretKey,
	}
}

func userIAMPbToInstance(pb *iam_pb.UserIAM) *UserIAM {
	if pb == nil {
		return nil
	}
	return &UserIAM{
		User:      User(pb.User),
		SecretKey: pb.SecretKey,
	}
}
