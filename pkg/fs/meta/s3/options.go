package s3

import "github.com/youngfs/youngfs/pkg/idmint"

type config struct {
	idGen idmint.Mint
}

type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (f optionFunc) apply(cfg *config) {
	f(cfg)
}

func WithIDGenerator(idGen idmint.Mint) Option {
	return optionFunc(func(cfg *config) {
		cfg.idGen = idGen
	})
}
