package engine

import (
	"context"
	"io"
)

type Engine interface {
	PutChunk(ctx context.Context, reader io.Reader, endpoints ...string) (string, error)
	GetChunk(ctx context.Context, id string) (io.ReadCloser, error)
	DeleteChunk(ctx context.Context, id string) error
	GetEndpoints(ctx context.Context) ([]string, error)
}
