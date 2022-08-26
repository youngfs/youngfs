package ec_calc

import (
	"bytes"
	"context"
	"github.com/go-playground/assert/v2"
	"icesfs/command/vars"
	"icesfs/ec/ec_store"
	"icesfs/kv/redis"
	"icesfs/storage_engine/seaweedfs"
	"icesfs/util"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestECCalc_ECReader(t *testing.T) {
	kvStore := redis.NewKvStore(vars.RedisSocket, vars.RedisPassword, vars.RedisDatabase)
	se := seaweedfs.NewStorageEngine(vars.SeaweedFSMaster, kvStore)

	ctx := context.Background()
	size := uint64(5 * 1024 * 1024)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 32; i++ {
		b := util.RandByte(size)
		err := error(nil)

		sz := make([]uint64, 3)
		fid := make([]string, 3)

		sz[0] = uint64(rand.Intn(int(size - 2)))
		sz[1] = uint64(rand.Intn(int(size - sz[0] - 1)))
		sz[2] = size - sz[0] - sz[1]

		fid[0], err = se.PutObject(ctx, sz[0], bytes.NewReader(b[:sz[0]]), "", true)
		assert.Equal(t, err, nil)

		fid[1], err = se.PutObject(ctx, sz[1], bytes.NewReader(b[sz[0]:sz[0]+sz[1]]), "", true)
		assert.Equal(t, err, nil)

		fid[2], err = se.PutObject(ctx, sz[2], bytes.NewReader(b[sz[0]+sz[1]:]), "", true)
		assert.Equal(t, err, nil)

		frags := make([]ec_store.Frag, 3)

		for i := 0; i < 3; i++ {
			frags[i] = ec_store.Frag{
				Fid:      fid[i],
				FileSize: sz[i],
			}
		}

		ecReader := NewECReadCloser(frags, se)
		cnt := uint64(0)
		for cnt < size {
			length := util.Max(uint64(rand.Intn(int(size-cnt))), 1)
			b2 := make([]byte, length)
			n, err := ecReader.Read(b2)
			assert.Equal(t, n, int(length))
			assert.Equal(t, err, nil)
			assert.Equal(t, util.BytesIsEqual(b2, b[cnt:cnt+length]), true)
			cnt += length
		}

		size2 := rand.Intn(1024)
		b3 := make([]byte, size2)
		b4 := make([]byte, size2)
		n, err := ecReader.Read(b4)
		assert.Equal(t, n, size2)
		assert.Equal(t, b4, b3)

		for i := 0; i < 3; i++ {
			err := se.DeleteObject(ctx, fid[i])
			assert.Equal(t, err, nil)
		}

		err = ecReader.Close()
		assert.Equal(t, err, nil)
	}

	for i := 0; i < 32; i++ {
		b := util.RandByte(size)
		err := error(nil)

		sz := make([]uint64, 3)
		fid := make([]string, 3)

		sz[0] = uint64(rand.Intn(int(size - 2)))
		sz[1] = uint64(rand.Intn(int(size - sz[0] - 1)))
		sz[2] = size - sz[0] - sz[1]

		fid[0], err = se.PutObject(ctx, sz[0], bytes.NewReader(b[:sz[0]]), "", true)
		assert.Equal(t, err, nil)

		fid[1] = ""

		fid[2], err = se.PutObject(ctx, sz[2], bytes.NewReader(b[sz[0]+sz[1]:]), "", true)
		assert.Equal(t, err, nil)

		frags := make([]ec_store.Frag, 3)

		for i := 0; i < 3; i++ {
			frags[i] = ec_store.Frag{
				Fid:      fid[i],
				FileSize: sz[i],
			}
		}

		for i := sz[0]; i < sz[0]+sz[1]; i++ {
			b[i] = 0
		}

		ecReader := NewECReadCloser(frags, se)
		cnt := uint64(0)
		for cnt < size {
			length := util.Max(uint64(rand.Intn(int(size-cnt))), 1)
			b2 := make([]byte, length)
			n, err := ecReader.Read(b2)
			assert.Equal(t, n, int(length))
			assert.Equal(t, err, nil)
			assert.Equal(t, util.BytesIsEqual(b2, b[cnt:cnt+length]), true)
			cnt += length
		}

		size2 := rand.Intn(1024)
		b3 := make([]byte, size2)
		b4 := make([]byte, size2)
		n, err := ecReader.Read(b4)
		assert.Equal(t, n, size2)
		assert.Equal(t, b4, b3)

		for i := 0; i < 3; i++ {
			if fid[i] != "" {
				err := se.DeleteObject(ctx, fid[i])
				assert.Equal(t, err, nil)
			}
		}

		err = ecReader.Close()
		assert.Equal(t, err, nil)
	}

	time.Sleep(3 * time.Second)
}

func TestECCalc_FilesReader(t *testing.T) {
	err := os.MkdirAll(ecFilePrefix, os.ModePerm)
	assert.Equal(t, err, nil)
	filesName := make([]string, 0)
	filesCnt := 0
	rand.Seed(time.Now().UnixNano())

	defer func() {
		for _, fileName := range filesName {
			err := os.Remove(fileName)
			assert.Equal(t, err, nil)
		}
		err := os.Remove(ecFilePrefix)
		assert.Equal(t, err, nil)
	}()

	size := 8 * 1024
	for i := 0; i < 32; i++ {
		b := util.RandByte(uint64(size))
		frags := make([]ec_store.Frag, 3)
		use := 0
		for j := 0; j < 3; j++ {
			fileName := ecFileKey("FilesReader", filesCnt)
			filesCnt++
			filesName = append(filesName, fileName)
			sz := 2*1024 + rand.Intn(1024)
			if j == 2 {
				sz = size - use
			}
			err := ioutil.WriteFile(fileName, b[use:use+sz], os.ModePerm)
			assert.Equal(t, err, nil)
			use += sz
			frags = append(frags, ec_store.Frag{
				Fid: fileName,
			})
		}

		emptyB := make([]byte, 8*1024)
		b = append(b, emptyB...)
		filesReader := NewFilesReader(frags)
		cnt := 0
		for cnt < size {
			length := 1024 + rand.Intn(512)
			filesReader.SetLimit(length)
			b2, err := ioutil.ReadAll(filesReader)
			assert.Equal(t, err, nil)
			assert.Equal(t, len(b2), length)
			assert.Equal(t, b2, b[cnt:cnt+length])
			cnt += length
		}
		err := filesReader.Release()
		assert.Equal(t, err, nil)
	}

	for i := 0; i < 32; i++ {
		b := util.RandByte(uint64(size))
		frags := make([]ec_store.Frag, 3)
		use := 0
		for j := 0; j < 3; j++ {
			fileName := ecFileKey("FilesReader", filesCnt)
			filesCnt++
			filesName = append(filesName, fileName)
			sz := 2*1024 + rand.Intn(1024)
			if j == 2 {
				sz = size - use
			}
			err := ioutil.WriteFile(fileName, b[use:use+sz], os.ModePerm)
			assert.Equal(t, err, nil)
			use += sz
			frags = append(frags, ec_store.Frag{
				Fid: fileName,
			})
		}

		emptyB := make([]byte, 8*1024)
		b = append(b, emptyB...)
		filesReader := NewFilesReader(frags)
		cnt := 0
		for cnt < (size >> 1) {
			length := 1024 + rand.Intn(512)
			filesReader.SetLimit(length)
			b2, err := ioutil.ReadAll(filesReader)
			assert.Equal(t, err, nil)
			assert.Equal(t, len(b2), length)
			assert.Equal(t, b2, b[cnt:cnt+length])
			cnt += length
		}
		err := filesReader.Release()
		assert.Equal(t, err, nil)
	}
}
