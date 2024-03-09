package seaweedfs

import (
	"github.com/oxtoacart/bpool"
	"github.com/youngfs/youngfs/pkg/log"
	"github.com/youngfs/youngfs/pkg/util"
	"github.com/youngfs/youngfs/pkg/util/gzippool"
	"sync"
)

type Engine struct {
	masterEndpoint string
	volumeIpMap    *sync.Map
	deletionQueue  *util.UnboundedQueue[string]
	gzipWriterPool *gzippool.GzipWriterPool
	bufferPool     *bpool.BufferPool
	endpoints      []string
	endpointsMutex *sync.RWMutex
	logger         log.Logger
}

func NewEngine(endpoint string, logger log.Logger) *Engine {
	e := &Engine{
		masterEndpoint: endpoint,
		volumeIpMap:    &sync.Map{},
		deletionQueue:  util.NewUnboundedQueue[string](),
		gzipWriterPool: gzippool.NewGzipWriterPool(),
		bufferPool:     bpool.NewBufferPool(128),
		endpoints:      make([]string, 0),
		endpointsMutex: &sync.RWMutex{},
		logger:         logger,
	}

	e.updateEndpoints()
	go e.scheduledUpdateEndpoints()
	return e
}
