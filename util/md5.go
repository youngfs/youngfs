package util

import (
	"crypto/md5"
)

func Md5ToBytes(b [md5.Size]byte) []byte {
	ret := make([]byte, md5.Size)

	for i, u := range b {
		ret[i] = u
	}

	return ret
}
