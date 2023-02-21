package randutil

import (
	"crypto/md5"
	"math/rand"
	"sync"
	"time"
)

var rp *RandPool
var defaultLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func init() {
	rp = NewRandPool()
}

type RandPool struct {
	sync.Pool
}

func NewRandPool() *RandPool {
	return &RandPool{
		Pool: sync.Pool{
			New: func() interface{} {
				w := rand.New(rand.NewSource(time.Now().Unix()))
				return w
			},
		},
	}
}

func (rp *RandPool) Read(p []byte) (n int, err error) {
	r, _ := rp.Get().(*rand.Rand)
	defer func() {
		rp.Put(r)
	}()
	return r.Read(p)
}

func (rp *RandPool) RandString(n int, letters []rune) string {
	r, _ := rp.Get().(*rand.Rand)
	defer func() {
		rp.Put(r)
	}()
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

func RandByte(len uint64) []byte {
	b := make([]byte, len)
	_, _ = rp.Read(b)
	return b
}

func RandMd5() []byte {
	// todo: update in go 1.20
	ret := make([]byte, md5.Size)
	sum := md5.Sum(RandByte(md5.BlockSize))
	copy(ret, sum[:])
	return ret
}

func RandString(n int, allowedChars ...[]rune) string {
	var letters []rune

	if len(allowedChars) == 0 {
		letters = defaultLetters
	} else {
		for _, chs := range allowedChars {
			letters = append(letters, chs...)
		}
	}

	return rp.RandString(n, letters)
}
