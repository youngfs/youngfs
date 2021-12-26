package kv

import (
	"github.com/go-redis/redis/v8"
)

type redisStore struct {
	client *redis.Client
}

var Client redisStore

func (store *redisStore) Initialize(hostPort, password string, database int) {
	store.client = redis.NewClient(
		&redis.Options{
			Addr:     hostPort,
			Password: password,
			DB:       database,
		},
	)
}
