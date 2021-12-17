package set_iam

import (
	"icesos/kv"
)

type User string

func (user User) UserIAMKey() string {
	return string(user) + userIAMKey
}

func (user User) SetIAMKey() string {
	return string(user) + setIAMKey
}

func (user User) IsOwnSet(set Set) bool {
	setIAM := SetIAM{
		User: user,
		Set:  set,
	}

	member, err := setIAM.EncodeProto()
	if err != nil {
		return false
	}

	ret, _ := kv.Client.SIsMember(user.SetIAMKey(), member)
	return ret
}

func (user User) Identify(sk string) bool {
	val, err := kv.Client.KvGet(user.UserIAMKey())
	if err != nil {
		return false
	}

	userIAM, err := DecodeUserIAMProto(val)
	if err != nil {
		return false
	}

	return userIAM.SecretKey == sk
}

func (user User) CreateUser(sk string) error {
	userIAM := UserIAM{
		User:      user,
		SecretKey: sk,
	}

	b, err := userIAM.EncodeProto()
	if err != nil {
		return err
	}

	return kv.Client.KvPut(user.UserIAMKey(), b)
}

func (user User) DeleteUser() error {
	_, err := kv.Client.KvDelete(user.UserIAMKey())
	if err != nil {
		return err
	}

	_, err = kv.Client.SDelete(user.SetIAMKey())
	return err
}
