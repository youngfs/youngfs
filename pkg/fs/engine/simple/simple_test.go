package simple

import (
	"github.com/stretchr/testify/suite"
	"github.com/youngfs/youngfs/pkg/fs/engine"
	"testing"
)

type simpleSuite struct {
	engine.EngineSuite
}

func TestSimpleStore(t *testing.T) {
	suite.Run(t, new(simpleSuite))
}

func (s *simpleSuite) SetupTest() {
	path := s.T().TempDir()
	e, err := New(path)
	s.Nil(err)
	s.Engine = e
}
