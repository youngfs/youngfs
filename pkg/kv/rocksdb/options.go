//go:build rocksdb
// +build rocksdb

package rocksdb

import "github.com/linxGnu/grocksdb"

type config struct {
	compactionFilter grocksdb.CompactionFilter
}

type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (f optionFunc) apply(cfg *config) {
	f(cfg)
}

func WithCompactionFilter(filter grocksdb.CompactionFilter) Option {
	return optionFunc(func(cfg *config) {
		cfg.compactionFilter = filter
	})
}
