package tikv

import "github.com/tikv/client-go/v2/txnkv"

type config struct {
	Enable1PC         bool
	EnablePessimistic bool
	options           []txnkv.ClientOpt
}

type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (f optionFunc) apply(cfg *config) {
	f(cfg)
}

func WithKeySpace(keyspace string) Option {
	return optionFunc(func(cfg *config) {
		cfg.options = append(cfg.options, txnkv.WithKeyspace(keyspace))
	})
}

func WithEnable1PC() Option {
	return optionFunc(func(cfg *config) {
		cfg.Enable1PC = true
	})
}

func WithEnablePessimistic() Option {
	return optionFunc(func(cfg *config) {
		cfg.EnablePessimistic = true
	})
}
