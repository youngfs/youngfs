package iam

import (
	"crypto/md5"
)

type userIAM struct {
	User                     // user_iam name
	SecretKey [md5.Size]byte // secret key md5
}

func (userIam *userIAM) key() string {
	return string(userIam.User) + userIAMKv
}
