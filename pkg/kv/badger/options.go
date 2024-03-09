package badger

import (
	"github.com/dgraph-io/badger/v4"
	"time"
)

type config struct {
	logger badger.Logger
	ttl    *time.Duration
}

type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (f optionFunc) apply(cfg *config) {
	f(cfg)
}

func WithLogger(logger badger.Logger) Option {
	return optionFunc(func(cfg *config) {
		cfg.logger = logger
	})
}

func WithTTL(ttl time.Duration) Option {
	return optionFunc(func(cfg *config) {
		cfg.ttl = &ttl
	})
}
