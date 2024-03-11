package master

import (
	"context"
	"github.com/youngfs/youngfs/pkg/chunk/pb/master_pb"
	"github.com/youngfs/youngfs/pkg/errors"
	"github.com/youngfs/youngfs/pkg/errors/ecode"
	"github.com/youngfs/youngfs/pkg/kv"
	"github.com/youngfs/youngfs/pkg/log"
	"github.com/youngfs/youngfs/pkg/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"strconv"
	"sync"
	"time"
)

type Server struct {
	master_pb.UnimplementedMasterServiceServer
	maxId      uint64
	volumeKv   kv.Store
	volumeLock *sync.RWMutex
	volume     map[string][]uint64
	endpoints  *sync.Map
	logger     log.Logger
}

func (s *Server) RegisterVolume(ctx context.Context, in *master_pb.RegisterVolumeRequest) (*master_pb.RegisterVolumeResponse, error) {
	s.volumeLock.Lock()
	defer s.volumeLock.Unlock()
	if len(in.Magic) == 0 {
		return nil, status.Errorf(codes.Code(ecode.ErrVolumeMagic), errors.ErrVolumeMagic.WithMessage("magic is empty").Error())
	}
	val, err := s.volumeKv.Get(ctx, in.Magic)
	if errors.Is(err, kv.ErrKeyNotFound) {
		s.maxId++
		err = s.volumeKv.Put(ctx, in.Magic, []byte(strconv.FormatUint(s.maxId, 10)))
		if err != nil {
			return nil, status.Errorf(codes.Code(ecode.ErrMaster), errors.ErrMaster.WarpErr(err).Error())
		}
		s.volume[in.Endpoint] = append(s.volume[in.Endpoint], s.maxId)
		return &master_pb.RegisterVolumeResponse{Id: s.maxId}, nil
	} else if err != nil {
		return nil, status.Errorf(codes.Code(ecode.ErrMaster), errors.ErrMaster.WarpErr(err).Error())
	} else {
		if in.Id == uint64(master_pb.ID_CreateId) {
			return nil, status.Errorf(codes.Code(ecode.ErrVolumeCreateConflict), errors.ErrVolumeCreateConflict.Error())
		}
		id, err := strconv.ParseUint(string(val), 10, 64)
		if err != nil {
			return nil, status.Errorf(codes.Code(ecode.ErrMaster), errors.ErrMaster.WithMessage("parse volume id failed").Error())
		}
		if id == in.Id {
			ids := append(s.volume[in.Endpoint], id)
			s.volume[in.Endpoint] = util.SortUnique(ids, func(i, j int) int {
				if ids[i] < ids[j] {
					return -1
				} else if ids[i] > ids[j] {
					return 1
				} else {
					return 0
				}
			})
			return &master_pb.RegisterVolumeResponse{Id: id}, nil
		} else {
			return nil, status.Errorf(codes.Code(ecode.ErrVolumeMagic), errors.ErrVolumeMagic.Error())
		}
	}
}

func (s *Server) SendHeartbeat(ctx context.Context, in *master_pb.HeartbeatRequest) (*emptypb.Empty, error) {
	s.endpoints.Store(in.Endpoint, time.Now())
	s.logger.Infof("received heartbeat from %s", in.Endpoint)
	return &emptypb.Empty{}, nil
}

func (s *Server) QueryEndpoints(ctx context.Context, _ *emptypb.Empty) (*master_pb.QueryResponse, error) {
	volumes := make([]*master_pb.VolumeInfo, 0)
	s.endpoints.Range(func(key, value interface{}) bool {
		endpoint, ok := key.(string)
		if !ok {
			return false
		}
		lastTime, ok := value.(time.Time)
		if !ok {
			return false
		}
		if time.Since(lastTime) < heartbeatTimeout {
			volumes = append(volumes, &master_pb.VolumeInfo{Endpoint: endpoint})
		}
		return true
	})
	s.volumeLock.RLock()
	defer s.volumeLock.RUnlock()
	for i, v := range volumes {
		volumes[i].Id = s.volume[v.Endpoint]
	}
	return &master_pb.QueryResponse{Volumes: volumes}, nil
}

func New(kv kv.Store, logger log.Logger) *Server {
	return &Server{
		volumeKv:   kv,
		endpoints:  &sync.Map{},
		logger:     logger,
		volumeLock: &sync.RWMutex{},
		volume:     make(map[string][]uint64),
	}
}

func (s *Server) Run(port int) error {
	it, err := s.volumeKv.NewIterator()
	if err != nil {
		return err
	}
	for ; it.Valid(); it.Next() {
		val := it.Value()
		id, err := strconv.ParseUint(string(val), 10, 64)
		if err != nil {
			it.Close()
			return err
		}
		if id > s.maxId {
			s.maxId = id
		}
	}
	it.Close()
	lis, err := net.Listen("tcp", strconv.Itoa(port))
	if err != nil {
		return err
	}
	srv := grpc.NewServer()
	master_pb.RegisterMasterServiceServer(srv, s)
	return srv.Serve(lis)
}
