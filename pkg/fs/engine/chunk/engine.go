package chunk

import (
	"context"
	"github.com/youngfs/youngfs/pkg/chunk/pb/master_pb"
	"github.com/youngfs/youngfs/pkg/chunk/pb/volume_pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
)

type Engine struct {
	master       string
	masterClient master_pb.MasterServiceClient
	volumeClient map[string]volume_pb.VolumeServiceClient
	endpoints    []string
	volumeMap    map[uint64]string
	mux          *sync.RWMutex
}

func New(master string) (*Engine, error) {
	conn, err := grpc.Dial(master, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	client := master_pb.NewMasterServiceClient(conn)
	e := &Engine{
		master:       master,
		masterClient: client,
		volumeClient: make(map[string]volume_pb.VolumeServiceClient),
		volumeMap:    make(map[uint64]string),
		mux:          &sync.RWMutex{},
	}
	err = e.updateEndpoints(context.Background())
	if err != nil {
		return nil, err
	}
	go e.scheduledUpdateEndpoints()
	return e, nil
}
