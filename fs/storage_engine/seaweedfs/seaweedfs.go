package seaweedfs

import (
	"sync"
	"youngfs/kv"
	"youngfs/util"
)

type StorageEngine struct {
	masterServer   string
	volumeIpMap    *sync.Map
	deletionQueue  *util.UnboundedQueue[string]
	kvStore        kv.KvSetStore
	gzipWriterPool *util.GzipWriterPool
}

func NewStorageEngine(masterServer string, KvStore kv.KvSetStore) *StorageEngine {
	se := &StorageEngine{
		masterServer:   masterServer,
		volumeIpMap:    &sync.Map{},
		deletionQueue:  util.NewUnboundedQueue[string](),
		kvStore:        KvStore,
		gzipWriterPool: util.NewGzipWriterPool(),
	}

	go se.loopProcessingDeletion()

	return se
}
