package badger

import (
	"github.com/dgraph-io/badger/v4"
	"time"
)

type Store struct {
	db  *badger.DB
	ttl *time.Duration
}

// New creates a new Store instance.
// path is a directory path where the database files will be stored.
func New(path string, opts ...Option) (*Store, error) {
	cfg := &config{}
	for _, opt := range opts {
		opt.apply(cfg)
	}

	dpOpts := badger.DefaultOptions(path)
	dpOpts = dpOpts.WithLogger(cfg.logger)
	db, err := badger.Open(dpOpts)
	if err != nil {
		return nil, err
	}
	return &Store{
		db:  db,
		ttl: cfg.ttl,
	}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}
