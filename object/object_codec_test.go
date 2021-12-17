package object

import (
	"github.com/go-playground/assert/v2"
	"icesos/util"
	rand2 "math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestObject_EnDecodeProto(t *testing.T) {
	b := util.RandByte(16)

	val := &Object{
		FullPath: "aa/bb/cc",
		Set:      "test",
		Time:     time.Now(),
		Mode:     os.ModeDir,
		Mime:     "",
		Md5:      b,
		FileSize: uint64(rand2.Int63()),
		VolumeId: uint64(rand2.Int63()),
		Fid:      strconv.Itoa(rand2.Int()),
	}
	val.Time = time.Unix(val.TimeUnix(), 0) // windows: precision to s
	b, err := val.EncodeProto()
	assert.Equal(t, err, nil)

	val2, err := DecodeObjectProto(b)
	assert.Equal(t, err, nil)
	assert.Equal(t, val2, val)
}
