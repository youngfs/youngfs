package util

import (
	"crypto/md5"
	"strconv"
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
	ret := ""

	for _, u := range b {
		str := strconv.FormatUint(uint64(u), 16)
		if len(str) == 1 {
			ret += "0"
		}
		ret += str
	}

	return ret
}
