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

//func (store *Redis3Store) InsertEntry(ctx context.Context, entry *object.Entry) error {
//
//	value, err := entry.EncodeProto()
//	if err != nil {
//		return fmt.Errorf("encoding %s %+v: %v", entry.FullPath, entry.Attribute, err)
//	}
//
//	err = store.KvPut(ctx, string(entry.FullPath), value)
//	if err != nil {
//		return fmt.Errorf("put %s: %v", entry.FullPath, err)
//	}
//
//	//dir, fileName := object.FullPath.DirAndName()
//
//	return nil
//}
