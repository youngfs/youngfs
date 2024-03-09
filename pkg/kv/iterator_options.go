package kv

type IteratorConfig struct {
	Prefix []byte
}

type IteratorOption interface {
	apply(*IteratorConfig)
}

type optionFunc func(*IteratorConfig)

func (f optionFunc) apply(cfg *IteratorConfig) {
	f(cfg)
}

func ApplyConfig(cfg *IteratorConfig, opts ...IteratorOption) {
	for _, opt := range opts {
		opt.apply(cfg)
	}
}

// WithPrefix returns an IteratorOption that sets the prefix for an iterator.
// This option configures the iterator to only access items that have the specified
// prefix. A prefix is a slice of bytes that the keys of items should start with
// in order to be included in the iteration.
func WithPrefix(prefix []byte) IteratorOption {
	return optionFunc(func(cfg *IteratorConfig) {
		cfg.Prefix = prefix
	})
}
