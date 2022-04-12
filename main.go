package main

import (
	"icesos/command/vars"
	"icesos/filer/vfs"
	"icesos/kv/redis"
	"icesos/routers"
	"icesos/server"
	"icesos/storage_engine/seaweedfs"
	"log"
)

func main() {
	kvStore := redis.NewKvStore(vars.RedisHostPost, vars.RedisPassword, vars.RedisDatabase)
	storageEngine := seaweedfs.NewStorageEngine(vars.MasterServer)
	filerStore := vfs.NewVFS(kvStore, storageEngine)
	server.Svr = server.NewServer(filerStore, storageEngine)

	//gin.SetMode(gin.ReleaseMode)
	r := routers.InitRouter()
	err := r.Run(":" + vars.Port)
	if err != nil {
		log.Println(err)
	}
	return
}
