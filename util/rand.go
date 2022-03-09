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

func RandMd5() []byte {
	return Md5ToBytes(md5.Sum(RandByte(md5.BlockSize)))
}
