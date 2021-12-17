package iam

import (
	"crypto/md5"
	"icesos/kv"
)

type User string

func (user User) userIAMKey() string {
	return string(user) + userIAMKv
}

func (user User) setIAMKey() string {
	return string(user) + setIAMKv
}

func (user User) IsOwnSet(set Set) bool {
	setIAM := setIAM{
		User: user,
		Set:  set,
	}

	member, err := setIAM.encodeProto()
	if err != nil {
		return false
	}

	ret, _ := kv.Client.SIsMember(setIAM.Key(), member)
	return ret
}

func (user User) AddSet(set Set) error {
	setIAM := setIAM{
		User: user,
		Set:  set,
	}

	member, err := setIAM.encodeProto()
	if err != nil {
		return err
	}

	return kv.Client.SAdd(setIAM.Key(), member)
}

func (user User) Identify(sk string) bool {
	val, err := kv.Client.KvGet(user.userIAMKey())
	if err != nil {
		return false
	}

	userIAM, err := decodeUserIAMProto(val)
	if err != nil {
		return false
	}

	return userIAM.SecretKey == md5.Sum([]byte(sk))
}

func (user User) CreateUser(sk string) error {
	userIAM := userIAM{
		User:      user,
		SecretKey: md5.Sum([]byte(sk)),
	}

	b, err := userIAM.encodeProto()
	if err != nil {
		return err
	}

	return kv.Client.KvPut(userIAM.key(), b)
}

func (user User) DeleteUser() (bool, error) {
	ret, err := kv.Client.KvDelete(user.userIAMKey())
	if err != nil || ret == false {
		return false, err
	}

	_, err = kv.Client.SDelete(user.setIAMKey())
	return true, err
}
