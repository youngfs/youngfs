package seaweedfs

import (
	"bytes"
	"context"
	"github.com/go-playground/assert/v2"
	"github.com/youngfs/youngfs/pkg/errors"
	"github.com/youngfs/youngfs/pkg/util"
	"github.com/youngfs/youngfs/pkg/util/randutil"
	"github.com/youngfs/youngfs/pkg/vars"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestSeaweedFS_PutObject(t *testing.T) {
	client := NewStorageEngine(vars.SeaweedFSMaster)
	size := uint64(5 * 1024)
	ctx := context.Background()
	rand.Seed(time.Now().UnixNano())

	b := randutil.RandByte(size)

	fid, err := client.PutObject(ctx, size, bytes.NewReader(b), "")
	assert.Equal(t, err, nil)

	url, err := client.getFidUrl(fid)
	assert.Equal(t, err, nil)

	buffer := &bytes.Buffer{}
	err = client.GetObject(ctx, fid, buffer)
	assert.Equal(t, err, nil)

	httpBody, err := io.ReadAll(buffer)
	assert.Equal(t, err, nil)
	assert.Equal(t, httpBody, b)

	err = client.DeleteObject(ctx, fid)
	assert.Equal(t, err, nil)

	time.Sleep(3 * time.Second)

	err = client.GetObject(ctx, fid, buffer)
	assert.Equal(t, errors.Is(err, errors.ErrObjectNotExist), true)

	resp, err := http.Get(url)
	defer func() {
		_ = resp.Body.Close()
	}()
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusNotFound)

	httpBody, err = io.ReadAll(buffer)
	assert.Equal(t, err, nil)
	assert.Equal(t, httpBody, []byte{})

	hosts, err := client.GetHosts(ctx)
	urls := make([]string, 0)

	lce := util.NewLimitedConcurrentExecutor(16)
	mutex := &sync.Mutex{}
	for i := 0; i < 128; i++ {
		lce.Execute(func() {
			b := randutil.RandByte(size)

			host := hosts[rand.Intn(len(hosts))]

			fid, err := client.PutObject(ctx, size, bytes.NewReader(b), host)
			assert.Equal(t, err, nil)

			volumeId, _, err := client.parseFid(fid)
			assert.Equal(t, err, nil)

			host2, err := client.getVolumeHost(volumeId)
			assert.Equal(t, err, nil)
			assert.Equal(t, host2, host)

			url, err := client.getFidUrl(fid)
			assert.Equal(t, err, nil)

			buffer := &bytes.Buffer{}
			err = client.GetObject(ctx, fid, buffer)
			assert.Equal(t, err, nil)

			httpBody, err := io.ReadAll(buffer)
			assert.Equal(t, err, nil)
			assert.Equal(t, httpBody, b)

			mutex.Lock()
			urls = append(urls, url)
			mutex.Unlock()

			err = client.DeleteObject(ctx, fid)
			assert.Equal(t, err, nil)
		})
	}

	lce.Wait()
	time.Sleep(5 * time.Second)

	for _, url_ := range urls {
		url := url_
		lce.Execute(func() {
			resp, err := http.Get(url)
			assert.Equal(t, err, nil)
			defer func() {
				_ = resp.Body.Close()
			}()
			assert.Equal(t, resp.StatusCode, http.StatusNotFound)

			httpBody, err = io.ReadAll(resp.Body)
			assert.Equal(t, err, nil)
			assert.Equal(t, httpBody, []byte{})
		})
	}
	lce.Wait()
}
