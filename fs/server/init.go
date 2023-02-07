package server

import (
	"youngfs/fs/ec/ec_calc"
	"youngfs/fs/ec/ec_server"
	"youngfs/fs/ec/ec_store"
	"youngfs/fs/filer/vfs"
	"youngfs/fs/id_generator/kv_generator"
	"youngfs/fs/storage_engine/seaweedfs"
	"youngfs/kv/redis"
	"youngfs/vars"
)

func InitServer() {
	kvStore := redis.NewKvStore(vars.RedisSocket, vars.RedisPassword, vars.RedisDatabase)
	storageEngine := seaweedfs.NewStorageEngine(vars.SeaweedFSMaster, kvStore)
	ecStore := ec_store.NewECStore(kvStore, storageEngine, kv_generator.NewKvGenerator(ecStoreIdGeneratorKey, kvStore))
	ecCalc := ec_calc.NewECCalc(ecStore, storageEngine)
	ecServer := ec_server.NewECServer(ecStore, ecCalc)
	filerStore := vfs.NewVFS(kvStore, storageEngine, ecServer)
	svr = NewServer(filerStore, storageEngine, ecServer)
}
