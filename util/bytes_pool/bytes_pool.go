package bytes_pool

import "sync"

var pools []*sync.Pool

func bitCount(size int) int {
	cnt := 0
	for ; size > minSize; cnt++ {
		size = (size + 1) >> 1
	}
	return cnt
}

func init() {
	// 1KB ~ 128MB
	pools = make([]*sync.Pool, bitCount(maxSize)+1)
	for i := 0; i < len(pools); i++ {
		size := minSize << i
		pools[i] = &sync.Pool{
			New: func() any {
				buffer := make([]byte, size)
				return &buffer
			},
		}
	}
}

func getBytesPool(size int) (*sync.Pool, bool) {
	index := bitCount(size)
	if index >= len(pools) {
		return nil, false
	}
	return pools[index], true
}

func Allocate(size int) []byte {
	if pool, ok := getBytesPool(size); ok {
		buf := *pool.Get().(*[]byte)
		return buf[:size]
	}
	return make([]byte, size)
}

func Free(buf []byte) {
	if pool, ok := getBytesPool(cap(buf)); ok {
		pool.Put(&buf)
	}
}
