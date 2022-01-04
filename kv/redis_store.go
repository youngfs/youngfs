package kv

import (
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

type redisStore struct {
	client  redis.UniversalClient
	redSync *redsync.Redsync
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
	store.redSync = redsync.New(goredis.NewPool(store.client))
}
