package storage_engine

type StorageEngine struct {
	masterServer string
	volumeIpMap  map[uint64]string
}

func NewStorageEngine(masterServer string) *StorageEngine {
	return &StorageEngine{
		masterServer: masterServer,
		volumeIpMap:  make(map[uint64]string),
	}
}
