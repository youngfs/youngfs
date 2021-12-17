package util

import "crypto/md5"

func Md5ToBytes(b [16]byte) []byte {
	ret := make([]byte, md5.Size)

	for i, u := range b {
		ret[i] = u
	}

	return ret
}

// use before check
func BytesToMd5(b []byte) [16]byte {
	var md5b [md5.Size]byte

	for i, u := range b {
		md5b[i] = u
	}

	return md5b
}
