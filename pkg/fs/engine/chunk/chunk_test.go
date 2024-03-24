package chunk

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"github.com/youngfs/youngfs/pkg/chunk/master"
	"github.com/youngfs/youngfs/pkg/chunk/volume"
	"github.com/youngfs/youngfs/pkg/chunk/volume/needle"
	"github.com/youngfs/youngfs/pkg/fs/engine"
	"github.com/youngfs/youngfs/pkg/kv"
	"github.com/youngfs/youngfs/pkg/kv/badger"
	"github.com/youngfs/youngfs/pkg/kv/bbolt"
	"github.com/youngfs/youngfs/pkg/log/zap"
	"io"
	"os"
	"path"
	"testing"
	"time"
)

const (
	masterPort = 9433
	slave1Port = 9434
	slave2Port = 9435
)

type chunkSuite struct {
	engine.EngineSuite
	closers []io.Closer
}

func TestChunkStore(t *testing.T) {
	suite.Run(t, new(chunkSuite))
}

func (s *chunkSuite) SetupTest() {
	dir := s.T().TempDir()
	logger := zap.New("chunkEngine", zap.WithDebug(), zap.WithLogWriter(io.Discard))

	mkv, err := bbolt.New(path.Join(dir, "master", "master.db"), []byte("master"))
	s.Nil(err)
	s.closers = append(s.closers, mkv)

	svr := master.New(mkv, logger)
	go func() {
		_ = svr.Run(masterPort)
	}()

	creator := needle.KvNeedleStoreCreator(func(path string) (kv.Store, error) {
		return badger.New(path)
	})
	err = os.MkdirAll(path.Join(dir, "volume1"), 0755)
	s.Nil(err)
	volume1 := volume.New(path.Join(dir, "volume1"), fmt.Sprintf("localhost:%d", masterPort), logger, creator)
	go func() {
		_ = volume1.Run(slave1Port)
	}()

	err = os.MkdirAll(path.Join(dir, "volume2"), 0755)
	s.Nil(err)
	volume2 := volume.New(path.Join(dir, "volume2"), fmt.Sprintf("localhost:%d", masterPort), logger, creator)
	go func() {
		_ = volume2.Run(slave2Port)
	}()

	time.Sleep(1 * time.Second)
	s.Engine, err = New(fmt.Sprintf("localhost:%d", masterPort))
	s.Nil(err)
}
