package bbolt

import (
	"go.etcd.io/bbolt"
	"os"
	pathutil "path"
	"strings"
)

type Store struct {
	db     *bbolt.DB
	bucket []byte
}

// New creates a new Store instance.
// path is a directory path or a file where the database files will be stored.
func New(path string, bucket []byte, opts ...Option) (*Store, error) {
	cfg := &config{}
	for _, opt := range opts {
		opt.apply(cfg)
	}

	pathInfo, err := os.Stat(path)
	name := "bbolt.db"
	if err != nil {
		if os.IsNotExist(err) {
			if !strings.HasSuffix(path, string(os.PathSeparator)) {
				name = pathutil.Base(path)
			}
			if err := os.MkdirAll(pathutil.Dir(path), 0755); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else if pathInfo.IsDir() {
		path = pathutil.Join(path, name)
	}

	dbOpt := bbolt.DefaultOptions
	dbOpt.NoSync = cfg.noSync
	// path need a file
	db, err := bbolt.Open(path, 0600, dbOpt)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(txn *bbolt.Tx) error {
		_, err := txn.CreateBucketIfNotExists(bucket)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &Store{
		db:     db,
		bucket: bucket,
	}, nil
}

func (s *Store) WithBucket(bucket []byte) (*Store, error) {
	err := s.db.Update(func(txn *bbolt.Tx) error {
		_, err := txn.CreateBucketIfNotExists(bucket)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &Store{
		db:     s.db,
		bucket: bucket,
	}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}
