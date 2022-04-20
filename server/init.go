package server

import (
	"icesos/command/vars"
	"icesos/filer/vfs"
	"icesos/kv/redis"
	"icesos/storage_engine/seaweedfs"
)

func InitServer() {
	kvStore := redis.NewKvStore(vars.RedisHostPost, vars.RedisPassword, vars.RedisDatabase)
	storageEngine := seaweedfs.NewStorageEngine(vars.MasterServer)
	filerStore := vfs.NewVFS(kvStore, storageEngine)
	svr = NewServer(filerStore, storageEngine)
}
