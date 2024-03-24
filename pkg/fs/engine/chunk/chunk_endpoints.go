package chunk

import (
	"context"
	"github.com/youngfs/youngfs/pkg/chunk/pb/volume_pb"
	"github.com/youngfs/youngfs/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"math/rand"
	"time"
)

func (e *Engine) GetEndpoints(ctx context.Context) ([]string, error) {
	e.mux.RLock()
	defer e.mux.RUnlock()
	ret := make([]string, len(e.endpoints))
	copy(ret, e.endpoints)
	return ret, nil
}

func (e *Engine) updateEndpoints(ctx context.Context) error {
	resp, err := e.masterClient.QueryEndpoints(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}
	e.mux.Lock()
	defer e.mux.Unlock()
	e.endpoints = e.endpoints[:0]
	clear(e.volumeMap)
	for _, volumes := range resp.Volumes {
		if e.volumeClient[volumes.Endpoint] == nil {
			conn, err := grpc.Dial(volumes.Endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
			if err != nil {
				return err
			}
			e.volumeClient[volumes.Endpoint] = volume_pb.NewVolumeServiceClient(conn)
		}
		e.endpoints = append(e.endpoints, volumes.Endpoint)
		for _, id := range volumes.Id {
			e.volumeMap[id] = volumes.Endpoint
		}
	}
	return nil
}

func (e *Engine) getVolumeClientFromId(id uint64) (volume_pb.VolumeServiceClient, error) {
	e.mux.RLock()
	defer e.mux.RUnlock()
	endpoint, ok := e.volumeMap[id]
	if !ok {
		return nil, errors.ErrVolumeNotFound
	}
	client, ok := e.volumeClient[endpoint]
	if !ok {
		return nil, errors.ErrEngineMaster
	}
	return client, nil
}

func (e *Engine) getVolumeClientFromEndpoints(endpoints ...string) (volume_pb.VolumeServiceClient, error) {
	e.mux.RLock()
	defer e.mux.RUnlock()
	endpoint := ""
	if len(endpoints) > 0 {
		endpoint = endpoints[rand.Intn(len(endpoints))]
	} else if len(e.endpoints) == 0 {
		return nil, errors.ErrEngineMaster
	} else {
		endpoint = e.endpoints[rand.Intn(len(e.endpoints))]
	}
	client, ok := e.volumeClient[endpoint]
	if !ok {
		return nil, errors.ErrEngineMaster
	}
	return client, nil
}

func (e *Engine) scheduledUpdateEndpoints() {
	ctx := context.Background()
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		_ = e.updateEndpoints(ctx)
	}
}
