package storage_engine

import "icesos/util"

type StorageEngine struct {
	masterServer  string
	volumeIpMap   map[uint64]string
	DeletionQueue *util.UnboundedQueue[string]
}

func NewStorageEngine(masterServer string) *StorageEngine {
	svr := &StorageEngine{
		masterServer:  masterServer,
		volumeIpMap:   make(map[uint64]string),
		DeletionQueue: util.NewUnboundedQueue[string](),
	}

	go svr.loopProcessingDeletion()

	return svr
}
