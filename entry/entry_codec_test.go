package entry

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
	val := &Entry{
		FullPath: "aa/bb/cc",
		Set:      "test",
		Time:     time.Now(),
		Mode:     os.ModeDir,
		Mime:     "",
		Md5:      util.RandMd5(),
		FileSize: uint64(rand2.Int63()),
		VolumeId: uint64(rand2.Int63()),
		Fid:      strconv.Itoa(rand2.Int()),
	}
	val.Time = time.Unix(val.Time.Unix(), 0) // windows: precision to s

	b, err := val.encodeProto()
	assert.Equal(t, err, nil)

	val2, err := decodeEntryProto(b)
	assert.Equal(t, err, nil)
	assert.Equal(t, val2, val)
}
