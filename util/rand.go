package util

import (
	"crypto/md5"
	"crypto/rand"
)

func RandByte(len uint64) []byte {
	b := make([]byte, len)
	_, _ = rand.Read(b)
	return b
}

func RandMd5() [16]byte {
	return md5.Sum(RandByte(md5.BlockSize))
}
