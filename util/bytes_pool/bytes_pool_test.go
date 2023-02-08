package bytes_pool

import (
	"bytes"
	"github.com/go-playground/assert/v2"
	"github.com/oxtoacart/bpool"
	"math/rand"
	"testing"
	"youngfs/util"
)

func TestBytesPool(t *testing.T) {
	for i := 0; i <= minSize; i++ {
		buf := Allocate(i)
		assert.Equal(t, len(buf), i)
		assert.Equal(t, cap(buf), minSize)
		assert.Equal(t, cap(buf), minSize<<bitCount(i))
		Free(buf)
	}
	for i := 0; i < 1024; i++ {
		size := rand.Intn(maxSize + (maxSize >> 2))
		buf := Allocate(size)
		assert.Equal(t, len(buf), size)
		if size <= maxSize {
			assert.Equal(t, cap(buf), minSize<<bitCount(size))
		} else {
			assert.Equal(t, cap(buf), size)
		}
		Free(buf)
	}
}

const testSize = 64 * 1024

func BenchmarkBytes(b *testing.B) {
	info := util.RandByte(testSize)
	for i := 0; i < b.N; i++ {
		var b []byte
		b = append(b, info...)
	}
}

func BenchmarkMakeBytes(b *testing.B) {
	info := util.RandByte(testSize)
	for i := 0; i < b.N; i++ {
		b := make([]byte, testSize)
		copy(b, info)
	}
}

func BenchmarkBuffer(b *testing.B) {
	info := util.RandByte(testSize)
	for i := 0; i < b.N; i++ {
		var b bytes.Buffer
		b.Write(info)
	}
}

func BenchmarkBufferAllocate(b *testing.B) {
	info := util.RandByte(testSize)
	for i := 0; i < b.N; i++ {
		b := bytes.NewBuffer(make([]byte, 0, testSize))
		b.Write(info)
	}
}

func BenchmarkBufferPool(b *testing.B) {
	info := util.RandByte(testSize)
	bufferPool := bpool.NewBufferPool(64)
	for i := 0; i < b.N; i++ {
		b := bufferPool.Get()
		b.Write(info)
		bufferPool.Put(b)
	}
}

func BenchmarkBytePool(b *testing.B) {
	info := util.RandByte(testSize)
	bytePool := bpool.NewBytePool(64, testSize)
	for i := 0; i < b.N; i++ {
		b := bytePool.Get()
		copy(b, info)
		bytePool.Put(b)
	}
}

func BenchmarkSizedBufferPool(b *testing.B) {
	info := util.RandByte(testSize)
	sizedBufferPool := bpool.NewSizedBufferPool(64, testSize)
	for i := 0; i < b.N; i++ {
		b := sizedBufferPool.Get()
		b.Write(info)
		sizedBufferPool.Put(b)
	}
}

func BenchmarkBytesPool(b *testing.B) {
	info := util.RandByte(testSize)
	for i := 0; i < b.N; i++ {
		b := Allocate(testSize)
		copy(b, info)
		Free(b)
	}
}
