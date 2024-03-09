//go:build rocksdb
// +build rocksdb

package rocksdb

import (
	"dmeta/pkg/kv"
	"github.com/stretchr/testify/suite"
	"testing"
)

type rocksdbSuite struct {
	kv.TransactionSuite
}

func TestRocksdbStore(t *testing.T) {
	suite.Run(t, new(rocksdbSuite))
}

func (s *rocksdbSuite) SetupTest() {
	path := s.T().TempDir()
	store, err := New(path)
	s.Nil(err)
	s.TransactionSuite.Store = store
	s.TransactionSuite.SetupTest()
}

func (s *rocksdbSuite) TearDownTest() {
	err := s.Store.Close()
	s.Nil(err)
}

type rocksdbTTLSuite struct {
	kv.TTLStoreSuite
}
