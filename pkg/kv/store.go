package kv

import (
	"context"
	"github.com/youngfs/youngfs/pkg/errors"
	"github.com/youngfs/youngfs/pkg/errors/ecode"
	"net/http"
)

var (
	ErrKeyNotFound = &errors.Error{Code: ecode.ErrKvNotFound, HTTPStatusCode: http.StatusContinue, Description: "Kv not found"}
)

type Store interface {
	Put(ctx context.Context, key []byte, val []byte) error
	Get(ctx context.Context, key []byte) ([]byte, error)
	Delete(ctx context.Context, key []byte) error
	NewIterator(opts ...IteratorOption) (Iterator, error)
	Close() error
}

type TransactionStore interface {
	Store
	NewTransaction() (Transaction, error)
}
