package seaweedfs

import (
	"github.com/oxtoacart/bpool"
	"sync"
	"youngfs/util"
	"youngfs/util/gzippool"
)

type StorageEngine struct {
	masterServer   string
	volumeIpMap    *sync.Map
	deletionQueue  *util.UnboundedQueue[string]
	gzipWriterPool *gzippool.GzipWriterPool
	bufferPool     *bpool.BufferPool
	hosts          []string
	hostsMutex     *sync.RWMutex
}

func NewStorageEngine(masterServer string) *StorageEngine {
	se := &StorageEngine{
		masterServer:   masterServer,
		volumeIpMap:    &sync.Map{},
		deletionQueue:  util.NewUnboundedQueue[string](),
		gzipWriterPool: gzippool.NewGzipWriterPool(),
		bufferPool:     bpool.NewBufferPool(128),
		hosts:          make([]string, 0),
		hostsMutex:     &sync.RWMutex{},
	}

	go se.loopProcessingDeletion()

	se.updateHosts()

	return se
}
