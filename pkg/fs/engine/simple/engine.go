package simple

import (
	"context"
	"fmt"
	"github.com/youngfs/youngfs/pkg/errors"
	"github.com/youngfs/youngfs/pkg/idmint"
	"github.com/youngfs/youngfs/pkg/idmint/uuid"
	"io"
	"os"
	"path/filepath"
)

type Engine struct {
	dir   string
	idGen idmint.Mint
}

func New(dir string, options ...Option) (*Engine, error) {
	cfg := &config{
		idGen: uuid.NewSimpleUUID(),
	}
	for _, opt := range options {
		opt.apply(cfg)
	}
	stat, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}
	if !stat.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", dir)
	}
	return &Engine{dir, cfg.idGen}, nil
}

func (e *Engine) PutChunk(ctx context.Context, reader io.Reader, endpoints ...string) (string, error) {
	id, err := e.idGen.String(ctx)
	if err != nil {
		return "", err
	}
	filePath := filepath.Join(e.dir, id)

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()

	_, err = io.Copy(file, reader)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (e *Engine) GetChunk(ctx context.Context, id string) (io.ReadCloser, error) {
	filePath := filepath.Join(e.dir, id)
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.ErrChunkNotExist
		}
		return nil, err
	}

	return file, nil
}

func (e *Engine) DeleteChunk(ctx context.Context, id string) error {
	filePath := filepath.Join(e.dir, id)
	return os.Remove(filePath)
}

func (e *Engine) GetEndpoints(ctx context.Context) ([]string, error) {
	return []string{"localhost"}, nil
}
