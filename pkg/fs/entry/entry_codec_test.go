package entry

import (
	"github.com/stretchr/testify/assert"
	"github.com/youngfs/youngfs/pkg/util/randutil"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestEntry_EnDecodeProto(t *testing.T) {
	val := &Entry{
		FullPath: "/aa/bb/cc",
		Bucket:   "test",
		Mtime:    time.Unix(time.Now().Unix(), 0),
		Ctime:    time.Unix(time.Now().Unix(), 0),
		Mode:     os.ModePerm,
		Mime:     "",
		Md5:      randutil.RandMd5(),
		FileSize: rand.Uint64(),
		Chunks: []*Chunk{
			{
				Offset:        rand.Uint64(),
				Size:          rand.Uint64(),
				Md5:           randutil.RandMd5(),
				IsReplication: rand.Int()%2 == 0,
				Frags: []*Frag{
					{
						Size:        rand.Uint64(),
						Id:          1,
						Md5:         randutil.RandMd5(),
						IsDataShard: rand.Int()%2 == 0,
						Fid:         randutil.RandString(16),
					},
					{
						Size:        rand.Uint64(),
						Id:          2,
						Md5:         randutil.RandMd5(),
						IsDataShard: rand.Int()%2 == 0,
						Fid:         randutil.RandString(16),
					},
				},
			},
			{
				Offset:        rand.Uint64(),
				Size:          rand.Uint64(),
				Md5:           randutil.RandMd5(),
				IsReplication: rand.Int()%2 == 0,
				Frags: []*Frag{
					{
						Size:        rand.Uint64(),
						Id:          1,
						Md5:         randutil.RandMd5(),
						IsDataShard: rand.Int()%2 == 0,
						Fid:         randutil.RandString(16),
					},
					{
						Size:        rand.Uint64(),
						Id:          2,
						Md5:         randutil.RandMd5(),
						IsDataShard: rand.Int()%2 == 0,
						Fid:         randutil.RandString(16),
					},
				},
			},
		},
	}

	b, err := val.EncodeProto()
	assert.Equal(t, err, nil)

	val2, err := DecodeEntryProto(b)
	assert.Equal(t, err, nil)
	assert.Equal(t, val2, val)

	val = &Entry{
		FullPath: "/aa/bb/cc",
		Bucket:   "test",
		Mtime:    time.Unix(time.Now().Unix(), 0),
		Ctime:    time.Unix(time.Now().Unix(), 0),
		Mode:     os.ModePerm,
		Mime:     "",
		FileSize: uint64(rand.Int63()),
		Chunks: []*Chunk{
			{
				Offset:        rand.Uint64(),
				Size:          rand.Uint64(),
				IsReplication: rand.Int()%2 == 0,
				Frags: []*Frag{
					{
						Size:        rand.Uint64(),
						Id:          1,
						IsDataShard: rand.Int()%2 == 0,
						Fid:         randutil.RandString(16),
					},
					{
						Size:        rand.Uint64(),
						Id:          2,
						IsDataShard: rand.Int()%2 == 0,
						Fid:         randutil.RandString(16),
					},
				},
			},
			{
				Offset:        rand.Uint64(),
				Size:          rand.Uint64(),
				IsReplication: rand.Int()%2 == 0,
				Frags: []*Frag{
					{
						Size:        rand.Uint64(),
						Id:          1,
						IsDataShard: rand.Int()%2 == 0,
						Fid:         randutil.RandString(16),
					},
					{
						Size:        rand.Uint64(),
						Id:          2,
						IsDataShard: rand.Int()%2 == 0,
						Fid:         randutil.RandString(16),
					},
				},
			},
		},
	}

	b, err = val.EncodeProto()
	assert.Equal(t, err, nil)

	val2, err = DecodeEntryProto(b)
	assert.Equal(t, err, nil)
	assert.Equal(t, val2, val)
}
