package ledisdb

import (
	lediscfg "github.com/ledisdb/ledisdb/config"
	"github.com/ledisdb/ledisdb/ledis"
	"github.com/youngfs/youngfs/errors"
	"github.com/youngfs/youngfs/log"
)

type KvStore struct {
	db *ledis.DB
}

func NewKvStore(dir, addr string) (*KvStore, error) {
	cfg := lediscfg.NewConfigDefault()
	cfg.Addr = addr
	cfg.DataDir = dir
	l, err := ledis.Open(cfg)
	if err != nil {
		log.Errorf("ledis init error :%s", err.Error())
		return nil, errors.ErrKvSever.WrapErr(err)
	}
	db, err := l.Select(0)
	if err != nil {
		log.Errorf("ledis init error :%s", err.Error())
		return nil, errors.ErrKvSever.WrapErr(err)
	}

	return &KvStore{
		db: db,
	}, nil
}
