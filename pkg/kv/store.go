package kv

import (
	"context"
	"github.com/youngfs/youngfs/pkg/errors"
	"net/http"
)

var (
	errKeyNotFound errors.Code = 1000
	ErrKeyNotFound             = &errors.Error{Code: errKeyNotFound, HTTPStatusCode: http.StatusContinue, Description: "Kv not found"}
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
