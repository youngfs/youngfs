package server

import (
	"github.com/youngfs/youngfs/pkg/fs/filer/vfs"
	"github.com/youngfs/youngfs/pkg/fs/storageengine/seaweedfs"
	"github.com/youngfs/youngfs/pkg/kv/redis"
	"github.com/youngfs/youngfs/pkg/vars"
)

func InitServer() {
	kvStore := redis.NewKvStore(vars.RedisSocket, vars.RedisPassword, vars.RedisDatabase)
	storageEngine := seaweedfs.NewStorageEngine(vars.SeaweedFSMaster)
	filerStore := vfs.NewVFS(kvStore, storageEngine)
	svr = NewServer(filerStore, storageEngine)
}
