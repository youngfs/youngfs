package kv

import (
	"context"
	"github.com/go-redis/redis/v8"
	"icesos/entry"
)

type Redis3Store struct {
	Client redis.UniversalClient
}

func (store *Redis3Store) Initialize(hostPort, password string, database int) {
	store.Client = redis.NewClient(
		&redis.Options{
			Addr:     hostPort,
			Password: password,
			DB:       database,
		},
	)
}

func (store *Redis3Store) InsertEntry(ctx context.Context, entry *entry.Entry) error {
	return nil
}
