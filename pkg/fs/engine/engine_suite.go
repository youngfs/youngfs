package engine

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/suite"
	"github.com/youngfs/youngfs/pkg/util/randutil"
	"io"
	"sync"
)

type EngineSuite struct {
	suite.Suite
	Engine Engine
}

func (s *EngineSuite) TestEngine() {
	s.NotNil(s.Engine)
	ctx := context.Background()
	const (
		size      = 1024
		chunkSize = (256 + 64) * 1024
	)

	type chunk struct {
		key  string
		body []byte
	}
	chunks := make([]chunk, size)
	wg := &sync.WaitGroup{}
	for i := range size {
		wg.Add(1)
		go func() {
			b := randutil.RandByte(chunkSize)
			id, err := s.Engine.PutChunk(ctx, bytes.NewReader(b))
			s.Nil(err)
			chunks[i] = chunk{id, b}
			wg.Done()
		}()
	}
	wg.Wait()

	for i := range size {
		wg.Add(1)
		go func() {
			reader, err := s.Engine.GetChunk(ctx, chunks[i].key)
			s.Nil(err)
			defer func() { _ = reader.Close() }()
			b, err := io.ReadAll(reader)
			s.Nil(err)
			s.Equal(chunks[i].body, b)
			wg.Done()
		}()
	}
	wg.Wait()

	for i := range size {
		wg.Add(1)
		go func() {
			err := s.Engine.DeleteChunk(ctx, chunks[i].key)
			s.Nil(err)
			wg.Done()
		}()
	}
	wg.Wait()

	for i := range size {
		wg.Add(1)
		go func() {
			reader, err := s.Engine.GetChunk(ctx, chunks[i].key)
			s.NotNil(err)
			s.Nil(reader)
			wg.Done()
		}()
	}
	wg.Wait()
}
