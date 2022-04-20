package entry

import (
	"context"
	"github.com/go-playground/assert/v2"
	"icesos/util"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestEntry_EnDecodeProto(t *testing.T) {
	ctx := context.Background()

	val := &Entry{
		FullPath: "/aa/bb/cc",
		Set:      "test",
		Mtime:    time.Unix(time.Now().Unix(), 0),
		Ctime:    time.Unix(time.Now().Unix(), 0),
		Mode:     os.ModePerm,
		Mime:     "",
		Md5:      util.RandMd5(),
		FileSize: uint64(rand.Int63()),
		Fid:      strconv.Itoa(rand.Int()),
		ECid:     strconv.FormatInt(rand.Int63(), 10),
	}

	b, err := val.EncodeProto(ctx)
	assert.Equal(t, err, nil)

	val2, err := DecodeEntryProto(ctx, b)
	assert.Equal(t, err, nil)
	assert.Equal(t, val2, val)

	val = &Entry{
		FullPath: "/aa/bb/cc",
		Set:      "test",
		Mtime:    time.Unix(time.Now().Unix(), 0),
		Ctime:    time.Unix(time.Now().Unix(), 0),
		Mode:     os.ModePerm,
		Mime:     "",
		FileSize: uint64(rand.Int63()),
		Fid:      strconv.Itoa(rand.Int()),
		ECid:     strconv.FormatInt(rand.Int63(), 10),
	}

	b, err = val.EncodeProto(ctx)
	assert.Equal(t, err, nil)

	val2, err = DecodeEntryProto(ctx, b)
	assert.Equal(t, err, nil)
	assert.Equal(t, val2, val)
}
