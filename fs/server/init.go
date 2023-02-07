package server

import (
	"youngfs/fs/filer/vfs"
	"youngfs/fs/storage_engine/seaweedfs"
	"youngfs/kv/redis"
	"youngfs/vars"
)

func InitServer() {
	kvStore := redis.NewKvStore(vars.RedisSocket, vars.RedisPassword, vars.RedisDatabase)
	storageEngine := seaweedfs.NewStorageEngine(vars.SeaweedFSMaster)
	filerStore := vfs.NewVFS(kvStore, storageEngine)
	svr = NewServer(filerStore, storageEngine)
}
