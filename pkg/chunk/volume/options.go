package volume

type config struct {
	localIP string
}

type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (f optionFunc) apply(cfg *config) {
	f(cfg)
}

func WithLocalIP(localIP string) Option {
	return optionFunc(func(cfg *config) {
		cfg.localIP = localIP
	})
}
