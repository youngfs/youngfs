package tests

import (
	"github.com/stretchr/testify/suite"
	"github.com/youngfs/youngfs/pkg/fs/engine/simple"
	"github.com/youngfs/youngfs/pkg/fs/handler"
	"github.com/youngfs/youngfs/pkg/fs/meta/s3"
	"github.com/youngfs/youngfs/pkg/fs/router"
	"github.com/youngfs/youngfs/pkg/fs/server"
	"github.com/youngfs/youngfs/pkg/kv/badger"
	"github.com/youngfs/youngfs/pkg/log/zap"
	"io"
	"net/http"
	"os"
	"path"
	"testing"
)

type handlerSuite struct {
	suite.Suite
	handler http.Handler
	closers []io.Closer
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(handlerSuite))
}

func (s *handlerSuite) SetupTest() {
	dir := s.T().TempDir()
	// logger
	logger := zap.New("youngfs", zap.WithDebug(), zap.WithLogWriter(io.Discard))
	// engine
	err := os.MkdirAll(path.Join(dir, "data"), 0755)
	s.Nil(err)
	engine, err := simple.New(path.Join(dir, "data"))
	s.Nil(err)
	// meta
	s3kv, err := badger.New(path.Join(dir, "meta", "s3kv"))
	s.Nil(err)
	s3cnkv, err := badger.New(path.Join(dir, "meta", "s3continuekv"))
	s.Nil(err)
	s.closers = append(s.closers, s3kv, s3cnkv)
	metaStore := s3.New(s3kv, s3cnkv)
	// server
	svr := server.NewServer(metaStore, engine)
	// handler
	h := handler.New(svr, logger)
	// router
	s.handler = router.New(h, router.WithMiddlewares(router.Logger(logger)))
}

func (s *handlerSuite) TearDownTest() {
	for _, c := range s.closers {
		err := c.Close()
		s.Nil(err)
	}
}
