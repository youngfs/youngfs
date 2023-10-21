package seaweedfs

import (
	"bytes"
	"context"
	"github.com/go-playground/assert/v2"
	"github.com/youngfs/youngfs/pkg/errors"
	"github.com/youngfs/youngfs/pkg/util/randutil"
	"github.com/youngfs/youngfs/pkg/vars"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestSeaweedFS_DeleteObject(t *testing.T) {
	client := NewStorageEngine(vars.SeaweedFSMaster)
	size := uint64(5 * 1024)
	ctx := context.Background()

	b := randutil.RandByte(size)

	fid, err := client.PutObject(ctx, size, bytes.NewReader(b))
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

	resp, err := http.Get(url)
	assert.Equal(t, err, nil)
	defer func() {
		_ = resp.Body.Close()
	}()
	assert.Equal(t, resp.StatusCode, http.StatusNotFound)

	err = client.GetObject(ctx, fid, buffer)
	assert.Equal(t, errors.Is(err, errors.ErrObjectNotExist), true)
}
