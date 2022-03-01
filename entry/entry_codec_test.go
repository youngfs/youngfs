package entry

import (
	"github.com/go-playground/assert/v2"
	"icesos/util"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestObject_EnDecodeProto(t *testing.T) {
	val := &Entry{
		FullPath: "/aa/bb/cc",
		Set:      "test",
		Ctime:    time.Unix(time.Now().Unix(), 0), // windows: precision to s
		Mode:     os.ModePerm,
		Mime:     "",
		Md5:      util.RandMd5(),
		FileSize: uint64(rand.Int63()),
		Fid:      strconv.Itoa(rand.Int()),
	}

	b, err := val.EncodeProto()
	assert.Equal(t, err, nil)

	val2, err := DecodeEntryProto(b)
	assert.Equal(t, err, nil)
	assert.Equal(t, val2, val)
}
