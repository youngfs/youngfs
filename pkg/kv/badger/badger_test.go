package badger

import (
	"github.com/stretchr/testify/suite"
	"github.com/youngfs/youngfs/pkg/kv"
	"testing"
	"time"
)

type badgerSuite struct {
	kv.TransactionSuite
}

func TestBadgerStore(t *testing.T) {
	suite.Run(t, new(badgerSuite))
}

func (s *badgerSuite) SetupTest() {
	path := s.T().TempDir()
	store, err := New(path)
	s.Nil(err)
	s.Store = store
	s.TransactionSuite.SetupTest()
}

func (s *badgerSuite) TearDownTest() {
	err := s.Store.Close()
	s.Nil(err)
}

type badgerTTLSuite struct {
	kv.TTLStoreSuite
}

func TestBadgerTTLStore(t *testing.T) {
	suite.Run(t, new(badgerTTLSuite))
}

func (s *badgerTTLSuite) SetupTest() {
	path := s.T().TempDir()
	s.TTL = 2 * time.Second
	store, err := New(path, WithTTL(s.TTL))
	s.Nil(err)
	s.Store = store
}

func (s *badgerTTLSuite) TearDownTest() {
	err := s.Store.Close()
	s.Nil(err)
}
