package directory

import (
	"github.com/go-playground/assert/v2"
	"os"
	"testing"
	"time"
)

func TestInode_EnDecodeProto(t *testing.T) {
	val := &Inode{
		FullPath: "/aa/bb/cc",
		Set:      "test",
		Time:     time.Now(),
		Mode:     os.ModeDir,
	}
	val.Time = time.Unix(val.Time.Unix(), 0) // windows: precision to s

	b, err := val.EncodeProto()
	assert.Equal(t, err, nil)

	val2, err := DecodeInodeProto(b)
	assert.Equal(t, err, nil)
	assert.Equal(t, val2, val)
}
