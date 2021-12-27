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
		Time:     time.Unix(time.Now().Unix(), 0), // windows: precision to s
		Mode:     os.ModeDir,
	}

	b, err := val.encodeProto()
	assert.Equal(t, err, nil)

	val2, err := decodeInodeProto(b)
	assert.Equal(t, err, nil)
	assert.Equal(t, val2, val)
}
