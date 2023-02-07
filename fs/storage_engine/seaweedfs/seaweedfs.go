package seaweedfs

import (
	"sync"
	"youngfs/util"
)

type StorageEngine struct {
	masterServer   string
	volumeIpMap    *sync.Map
	deletionQueue  *util.UnboundedQueue[string]
	gzipWriterPool *util.GzipWriterPool
	hosts          []string
	hostsMutex     *sync.RWMutex
}

func NewStorageEngine(masterServer string) *StorageEngine {
	se := &StorageEngine{
		masterServer:   masterServer,
		volumeIpMap:    &sync.Map{},
		deletionQueue:  util.NewUnboundedQueue[string](),
		gzipWriterPool: util.NewGzipWriterPool(),
		hosts:          make([]string, 0),
		hostsMutex:     &sync.RWMutex{},
	}

	go se.loopProcessingDeletion()

	se.updateHosts()

	return se
}
