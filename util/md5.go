package util

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5ToBytes(b [md5.Size]byte) []byte {
	ret := make([]byte, md5.Size)

	for i, u := range b {
		ret[i] = u
	}

	return ret
}

// use before check md5 size
func BytesToMd5(b []byte) [md5.Size]byte {
	var md5b [md5.Size]byte

	for i, u := range b {
		md5b[i] = u
	}

	return md5b
}

func Md5ToStr(b [md5.Size]byte) string {
	return hex.EncodeToString(Md5ToBytes(b))
}

func Md5IsEmpty(b [md5.Size]byte) bool {
	for _, u := range b {
		if u != 0 {
			return false
		}
	}
	return true
}
