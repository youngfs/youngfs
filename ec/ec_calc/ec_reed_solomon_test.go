package ec_calc

import (
	"github.com/go-playground/assert/v2"
	"github.com/klauspost/reedsolomon"
	"icesfs/util"
	"math/rand"
	"testing"
	"time"
)

func TestECCalc_ReedSolomonPackage(t *testing.T) {
	enc, err := reedsolomon.New(4, 2)
	assert.Equal(t, err, nil)
	size := uint64(64 * 1024 * 1024)
	emptyBytes := make([]byte, size)
	rand.Seed(time.Now().UnixNano())

	data := make([][]byte, 6)
	for i := 0; i < 4; i++ {
		data[i] = util.RandByte(size)
	}
	for i := 4; i < 6; i++ {
		data[i] = make([]byte, size)
		assert.Equal(t, data[i], emptyBytes)
	}

	err = enc.Encode(data)
	assert.Equal(t, err, nil)
	for i := 4; i < 6; i++ {
		assert.Equal(t, util.BytesIsEqual(data[i], emptyBytes), false)
	}

	for i := 0; i < 32; i++ {
		x := make([]int, 2)
		for j := range x {
			x[j] = rand.Intn(6)
		}
		datax := make([][]byte, 2)

		for j := range x {
			datax[j] = make([]byte, size)
			copy(datax[j], data[x[j]])
		}

		for j := range x {
			data[x[j]] = nil
		}

		err := enc.Reconstruct(data)
		assert.Equal(t, err, nil)

		for j := range x {
			assert.Equal(t, datax[j], data[x[j]])
		}
	}

	enc, err = reedsolomon.New(3, 3)
	assert.Equal(t, err, nil)

	for i := 0; i < 3; i++ {
		data[i] = util.RandByte(size)
	}
	for i := 3; i < 6; i++ {
		data[i] = make([]byte, size)
		assert.Equal(t, data[i], emptyBytes)
	}

	err = enc.Encode(data)
	assert.Equal(t, err, nil)
	for i := 3; i < 6; i++ {
		assert.Equal(t, util.BytesIsEqual(data[i], emptyBytes), false)
	}

	for i := 0; i < 32; i++ {
		x := make([]int, 3)
		for j := range x {
			x[j] = rand.Intn(6)
		}
		datax := make([][]byte, 3)

		for j := range x {
			datax[j] = make([]byte, size)
			copy(datax[j], data[x[j]])
		}

		for j := range x {
			data[x[j]] = nil
		}

		err := enc.Reconstruct(data)
		assert.Equal(t, err, nil)

		for j := range x {
			assert.Equal(t, datax[j], data[x[j]])
		}
	}

	enc, err = reedsolomon.New(2, 4)
	assert.Equal(t, err, nil)

	for i := 0; i < 2; i++ {
		data[i] = util.RandByte(size)
	}
	for i := 2; i < 6; i++ {
		data[i] = make([]byte, size)
		assert.Equal(t, data[i], emptyBytes)
	}

	err = enc.Encode(data)
	assert.Equal(t, err, nil)
	for i := 2; i < 6; i++ {
		assert.Equal(t, util.BytesIsEqual(data[i], emptyBytes), false)
	}

	for i := 0; i < 32; i++ {
		x := make([]int, 4)
		for j := range x {
			x[j] = rand.Intn(6)
		}
		datax := make([][]byte, 4)

		for j := range x {
			datax[j] = make([]byte, size)
			copy(datax[j], data[x[j]])
		}

		for j := range x {
			data[x[j]] = nil
		}

		err := enc.Reconstruct(data)
		assert.Equal(t, err, nil)

		for j := range x {
			assert.Equal(t, datax[j], data[x[j]])
		}
	}
}
