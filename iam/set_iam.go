package set_iam

import (
	"icesos/iam/iam_pb"
)

type SetIAM struct {
	User // user_iam name
	Set  // set_iam name
}

func (setIam *SetIAM) Key() string {
	return string(setIam.User) + setIAMKey
}

func (setIam *SetIAM) toPb() *iam_pb.SetIAM {
	if setIam == nil {
		return nil
	}
	return &iam_pb.SetIAM{
		User: string(setIam.User),
		Set:  string(setIam.Set),
	}
}

func setIAMPbToInstance(pb *iam_pb.SetIAM) *SetIAM {
	if pb == nil {
		return nil
	}
	return &SetIAM{
		User: User(pb.User),
		Set:  Set(pb.Set),
	}
}
