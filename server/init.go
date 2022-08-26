package server

import (
	"icesfs/command/vars"
	"icesfs/ec/ec_calc"
	"icesfs/ec/ec_server"
	"icesfs/ec/ec_store"
	"icesfs/filer/vfs"
	"icesfs/kv/redis"
	"icesfs/storage_engine/seaweedfs"
)

func InitServer() {
	kvStore := redis.NewKvStore(vars.RedisSocket, vars.RedisPassword, vars.RedisDatabase)
	storageEngine := seaweedfs.NewStorageEngine(vars.SeaweedFSMaster, kvStore)
	ecStore := ec_store.NewEC(kvStore, storageEngine)
	ecCalc := ec_calc.NewECCalc(ecStore, storageEngine)
	ecServer := ec_server.NewECServer(ecStore, ecCalc)
	filerStore := vfs.NewVFS(kvStore, storageEngine, ecServer)
	svr = NewServer(filerStore, storageEngine, ecServer)
}
