package bbolt

import (
	"github.com/stretchr/testify/suite"
	"github.com/youngfs/youngfs/pkg/kv"
	"testing"
)

type bboltSuite struct {
	kv.StoreSuite // bboltDB does not support multiple read and write transactions simultaneously.
}

func TestBboltStore(t *testing.T) {
	suite.Run(t, new(bboltSuite))
}

func (s *bboltSuite) SetupTest() {
	path := s.T().TempDir()
	store, err := New(path, []byte("test"))
	s.Nil(err)
	s.Store = store
}

func (s *bboltSuite) TearDownTest() {
	err := s.Store.Close()
	s.Nil(err)
}
