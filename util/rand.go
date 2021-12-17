package util

import "crypto/rand"

func RandByte(len int) []byte {
	b := make([]byte, len)
	_, _ = rand.Read(b)
	return b
}
