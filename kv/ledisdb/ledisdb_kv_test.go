package ledisdb

import (
	"context"
	"github.com/go-playground/assert/v2"
	"github.com/youngfs/youngfs/errors"
	"github.com/youngfs/youngfs/util/randutil"
	"os"
	"testing"
)

func TestLedis_Kv(t *testing.T) {
	client, err := NewKvStore(".kv", "")
	assert.Equal(t, err, nil)
	defer func() {
		_ = os.RemoveAll(".kv")
	}()

	key := "test_ledis_kv"
	ctx := context.Background()

	b, err := client.KvGet(ctx, key)
	assert.Equal(t, errors.IsKvNotFound(err), true)
	assert.Equal(t, b, nil)

	ret, err := client.KvDelete(ctx, key)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, true)

	b = randutil.RandByte(1024)
	err = client.KvPut(ctx, key, b)
	assert.Equal(t, err, nil)

	b2, err := client.KvGet(ctx, key)
	assert.Equal(t, err, nil)
	assert.Equal(t, b2, b)

	b = randutil.RandByte(512)
	err = client.KvPut(ctx, key, b)
	assert.Equal(t, err, nil)

	b2, err = client.KvGet(ctx, key)
	assert.Equal(t, err, nil)
	assert.Equal(t, b2, b)

	ret, err = client.KvDelete(ctx, key)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, true)

	b, err = client.KvGet(ctx, key)
	assert.Equal(t, errors.IsKvNotFound(err), true)
	assert.Equal(t, b, nil)

	ret, err = client.KvDelete(ctx, key)
	assert.Equal(t, err, nil)
	assert.Equal(t, ret, true)
}
