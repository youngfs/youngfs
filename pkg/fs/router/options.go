package router

import "github.com/gin-gonic/gin"

type config struct {
	debug       bool
	middlewares []gin.HandlerFunc
}

type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (f optionFunc) apply(cfg *config) {
	f(cfg)
}

func WithDebug(debug bool) Option {
	return optionFunc(func(cfg *config) {
		cfg.debug = debug
	})
}

func WithMiddlewares(middlewares ...gin.HandlerFunc) Option {
	return optionFunc(func(cfg *config) {
		cfg.middlewares = append(cfg.middlewares, middlewares...)
	})
}
