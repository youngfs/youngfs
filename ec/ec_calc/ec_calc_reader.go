package ec_calc

import (
	"context"
	"icesos/ec/ec_store"
	"icesos/storage_engine"
	"icesos/util"
	"io"
	"math"
	"net/http"
	"os"
	"sync"
)

type EmptyReadCloser struct {
	cnt    int
	closed bool
	locker sync.Mutex
}

func (reader *EmptyReadCloser) Read(p []byte) (int, error) {
	reader.locker.Lock()
	defer reader.locker.Unlock()

	if reader.closed || reader.cnt == 0 {
		return 0, io.EOF
	}

	mn := util.Min(reader.cnt, len(p))
	for i := 0; i < mn; i++ {
		p[i] = 0
	}

	reader.cnt -= mn

	if reader.cnt == 0 {
		return mn, io.EOF
	} else {
		return mn, nil
	}
}

func (reader *EmptyReadCloser) Close() error {
	reader.locker.Lock()
	defer reader.locker.Unlock()

	reader.closed = true
	return nil
}

func NewEmptyReadCloser(n int) *EmptyReadCloser {
	return &EmptyReadCloser{
		cnt:    n,
		closed: false,
	}
}

type ECReadCloser struct {
	frags         []ec_store.Frag
	p             int
	storageEngine storage_engine.StorageEngine
	file          io.ReadCloser
	closed        bool
	locker        sync.Mutex
}

func (reader *ECReadCloser) Read(p []byte) (int, error) {
	reader.locker.Lock()
	defer reader.locker.Unlock()

	if reader.closed {
		return 0, io.EOF
	}

	cnt := 0
first: // not change cnt

	for reader.file != nil {
		n, err := reader.file.Read(p[cnt:])
		cnt += n

		if cnt == len(p) {
			return cnt, nil
		}

		if err == io.EOF {
			_ = reader.file.Close()
			reader.file = nil
		}
	}

	if reader.p >= len(reader.frags) {
		reader.file = NewEmptyReadCloser(math.MaxInt)
		goto first
	}

	url, err := reader.storageEngine.GetFidUrl(context.Background(), reader.frags[reader.p].Fid)
	if err != nil {
		reader.file = NewEmptyReadCloser(int(reader.frags[reader.p].FileSize))
	} else {
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK || util.GetContentLength(resp.Header) != reader.frags[reader.p].FileSize {
			reader.file = NewEmptyReadCloser(int(reader.frags[reader.p].FileSize))
		} else {
			reader.file = resp.Body
		}
	}

	reader.p++
	goto first
}

func (reader *ECReadCloser) Close() error {
	reader.locker.Lock()
	defer reader.locker.Unlock()

	if reader.file != nil {
		err := reader.file.Close()
		if err != nil {
			return err
		}
	}

	reader.closed = true
	return nil
}

func NewECReadCloser(frags []ec_store.Frag, se storage_engine.StorageEngine) *ECReadCloser {
	fragsCopy := make([]ec_store.Frag, len(frags))
	copy(fragsCopy, frags)

	return &ECReadCloser{
		frags:         fragsCopy,
		p:             0,
		storageEngine: se,
		file:          nil,
	}
}

type FilesReader struct {
	frags  []ec_store.Frag
	p      int
	file   io.ReadCloser
	locker sync.Mutex
	limit  int
}

func (reader *FilesReader) Read(p []byte) (int, error) {
	reader.locker.Lock()
	defer reader.locker.Unlock()

	if reader.limit == 0 {
		return 0, io.EOF
	}
	cnt := 0
first: // not change cnt

	for reader.file != nil {
		n, err := reader.file.Read(p[cnt:util.Min(len(p), cnt+reader.limit)])
		cnt += n
		reader.limit -= n

		if reader.limit == 0 {
			return cnt, io.EOF
		}

		if cnt == len(p) {
			return cnt, nil
		}

		if err == io.EOF {
			_ = reader.file.Close()
			reader.file = nil
		}
	}

	if reader.p >= len(reader.frags) {
		reader.file = NewEmptyReadCloser(math.MaxInt)
		goto first
	}

	err := error(nil)
	reader.file, err = os.Open(reader.frags[reader.p].Fid)
	if err != nil {
		reader.file = nil
	}

	reader.p++
	goto first
}

func (reader *FilesReader) Release() error {
	reader.locker.Lock()
	defer reader.locker.Unlock()

	if reader.file != nil {
		err := reader.file.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (reader *FilesReader) SetLimit(limit int) {
	reader.locker.Lock()
	defer reader.locker.Unlock()

	reader.limit = limit
}

func NewFilesReader(frags []ec_store.Frag) *FilesReader {
	fragsCopy := make([]ec_store.Frag, len(frags))
	copy(fragsCopy, frags)

	return &FilesReader{
		frags: fragsCopy,
		p:     0,
		file:  nil,
		limit: 0,
	}
}
