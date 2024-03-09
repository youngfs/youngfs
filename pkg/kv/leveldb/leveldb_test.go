package leveldb

import (
	"github.com/stretchr/testify/suite"
	"github.com/youngfs/youngfs/pkg/kv"
	"testing"
)

type levelDBSuite struct {
	kv.StoreSuite // levelDB does not support multiple read and write transactions simultaneously.
}

func TestBadgerStore(t *testing.T) {
	suite.Run(t, new(levelDBSuite))
}

func (s *levelDBSuite) SetupTest() {
	path := s.T().TempDir()
	store, err := New(path)
	s.Nil(err)
	s.Store = store
}

func (s *levelDBSuite) TearDownTest() {
	err := s.Store.Close()
	s.Nil(err)
}
