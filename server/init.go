package server

import (
	"icesos/command/vars"
	"icesos/ec/ec_calc"
	"icesos/ec/ec_server"
	"icesos/ec/ec_store"
	"icesos/filer/vfs"
	"icesos/kv/redis"
	"icesos/storage_engine/seaweedfs"
)

func InitServer() {
	kvStore := redis.NewKvStore(vars.RedisHostPost, vars.RedisPassword, vars.RedisDatabase)
	storageEngine := seaweedfs.NewStorageEngine(vars.MasterServer, kvStore)
	ecStore := ec_store.NewEC(kvStore, storageEngine)
	ecCalc := ec_calc.NewECCalc(ecStore, storageEngine)
	ecServer := ec_server.NewECServer(ecStore, ecCalc)
	filerStore := vfs.NewVFS(kvStore, storageEngine, ecServer)
	svr = NewServer(filerStore, storageEngine, ecServer)
}
