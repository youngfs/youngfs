package bbolt

type config struct {
	noSync bool
}

type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (f optionFunc) apply(cfg *config) {
	f(cfg)
}

func WithNoSync() Option {
	return optionFunc(func(cfg *config) {
		cfg.noSync = true
	})
}
