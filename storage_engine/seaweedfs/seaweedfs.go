package seaweedfs

import (
	"icesos/kv"
	"icesos/util"
)

type StorageEngine struct {
	masterServer  string
	volumeIpMap   map[uint64]string
	deletionQueue *util.UnboundedQueue[string]
	kvStore       kv.KvStore
}

func NewStorageEngine(masterServer string, KvStore kv.KvStore) *StorageEngine {
	se := &StorageEngine{
		masterServer:  masterServer,
		volumeIpMap:   make(map[uint64]string),
		deletionQueue: util.NewUnboundedQueue[string](),
		kvStore:       KvStore,
	}

	go se.loopProcessingDeletion()

	return se
}
