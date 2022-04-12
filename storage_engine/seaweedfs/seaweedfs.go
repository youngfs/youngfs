package seaweedfs

import "icesos/util"

type StorageEngine struct {
	masterServer  string
	volumeIpMap   map[uint64]string
	deletionQueue *util.UnboundedQueue[string]
}

func NewStorageEngine(masterServer string) *StorageEngine {
	se := &StorageEngine{
		masterServer:  masterServer,
		volumeIpMap:   make(map[uint64]string),
		deletionQueue: util.NewUnboundedQueue[string](),
	}

	go se.loopProcessingDeletion()

	return se
}
