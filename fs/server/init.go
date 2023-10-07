package server

import (
	"github.com/youngfs/youngfs/fs/filer/vfs"
	"github.com/youngfs/youngfs/fs/storageengine/seaweedfs"
	"github.com/youngfs/youngfs/kv/redis"
	"github.com/youngfs/youngfs/vars"
)

func InitServer() {
	kvStore := redis.NewKvStore(vars.RedisSocket, vars.RedisPassword, vars.RedisDatabase)
	storageEngine := seaweedfs.NewStorageEngine(vars.SeaweedFSMaster)
	filerStore := vfs.NewVFS(kvStore, storageEngine)
	svr = NewServer(filerStore, storageEngine)
}
