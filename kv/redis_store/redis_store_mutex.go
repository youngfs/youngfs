package redis_store

import "github.com/go-redsync/redsync/v4"

func (store *RedisStore) NewMutex(name string, options ...redsync.Option) *redsync.Mutex {
	return store.redSync.NewMutex(name, options...)
}
