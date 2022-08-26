package util

import (
	"compress/gzip"
	"icesfs/errors"
	"icesfs/log"
	"io"
	"sync"
)

type GzipWriterPool struct {
	sync.Pool
}

func NewGzipWriterPool() *GzipWriterPool {
	return &GzipWriterPool{
		Pool: sync.Pool{
			New: func() interface{} {
				w, _ := gzip.NewWriterLevel(nil, gzip.BestSpeed)
				return w
			},
		},
	}
}

func (pool *GzipWriterPool) GzipStream(w io.Writer, r io.Reader) (int64, error) {
	gw, ok := pool.Get().(*gzip.Writer)
	if !ok {
		log.Errorw("gzip: new writer error")
		return 0, errors.GetAPIErr(errors.ErrServer)
	}
	gw.Reset(w)
	defer func() {
		_ = gw.Close()
		pool.Put(gw)
	}()
	return io.Copy(gw, r)
}
