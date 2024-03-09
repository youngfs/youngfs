package simple

import "github.com/youngfs/youngfs/pkg/idgenerator"

type config struct {
	idGen idgenerator.IDGenerator
}

type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (f optionFunc) apply(cfg *config) {
	f(cfg)
}

func WithIDGenerator(idGen idgenerator.IDGenerator) Option {
	return optionFunc(func(cfg *config) {
		cfg.idGen = idGen
	})
}
