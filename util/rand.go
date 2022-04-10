package util

import (
	"crypto/md5"
	"math/rand"
	"time"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix()))
}

func RandByte(len uint64) []byte {
	b := make([]byte, len)
	_, _ = r.Read(b)
	return b
}

func RandMd5() []byte {
	return Md5ToBytes(md5.Sum(RandByte(md5.BlockSize)))
}

var defaultLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandString(n int, allowedChars ...[]rune) string {
	var letters []rune

	if len(allowedChars) == 0 {
		letters = defaultLetters
	} else {
		for _, chs := range allowedChars {
			letters = append(letters, chs...)
		}
	}

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}

	return string(b)
}
