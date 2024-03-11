package volume

import (
	"context"
	"fmt"
	"github.com/Workiva/go-datastructures/queue"
	"github.com/youngfs/youngfs/pkg/chunk/pb/master_pb"
	"github.com/youngfs/youngfs/pkg/chunk/pb/volume_pb"
	"github.com/youngfs/youngfs/pkg/chunk/volume/needle"
	"github.com/youngfs/youngfs/pkg/errors"
	"github.com/youngfs/youngfs/pkg/errors/ecode"
	"github.com/youngfs/youngfs/pkg/log"
	"github.com/youngfs/youngfs/pkg/util/mem"
	"github.com/youngfs/youngfs/pkg/util/netutil"
	"github.com/youngfs/youngfs/pkg/util/randutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"net"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Server struct {
	volume_pb.UnimplementedVolumeServiceServer
	dir           string
	master        string
	client        master_pb.MasterServiceClient
	localEndpoint string
	logger        log.Logger
	queue         *queue.RingBuffer
	creator       needle.StoreCreator
	volumeCount   *NumberFile
	volumeMap     *sync.Map
}

func New(dir, master string, logger log.Logger, creator needle.StoreCreator) *Server {
	return &Server{
		dir:         dir,
		master:      master,
		logger:      logger,
		queue:       queue.NewRingBuffer(maxWritableVolume),
		creator:     creator,
		volumeCount: NewNumberFile(path.Join(dir, numFileName)),
	}
}

func (s *Server) Run(port int) error {
	ip, err := netutil.LocalIP()
	if err != nil {
		return err
	}
	s.localEndpoint = fmt.Sprintf("%s:%d", ip[0], port)

	conn, err := grpc.Dial(s.master, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()
	s.client = master_pb.NewMasterServiceClient(conn)
	go s.SendHeartbeat()
	lis, err := net.Listen("tcp", strconv.Itoa(port))
	if err != nil {
		return err
	}
	srv := grpc.NewServer()
	volume_pb.RegisterVolumeServiceServer(srv, s)
	return srv.Serve(lis)
}

func (s *Server) SendHeartbeat() {
	ctx := context.Background()
	for range time.NewTicker(heartbeatTick).C {
		_, err := s.client.SendHeartbeat(ctx, &master_pb.HeartbeatRequest{
			Endpoint: s.localEndpoint,
		})
		if err != nil {
			s.logger.Errorf("send heartbeat failed: %v", err)
		} else {
			s.logger.Infof("send heartbeat success")
		}
	}
}

func (s *Server) PutChunk(stream volume_pb.VolumeService_PutChunkServer) error {
	reader, writer := io.Pipe()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	var id string
	var err error
	go func() {
		defer wg.Done()
		id, err = s.putChunk(reader)
	}()
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			if err := writer.Close(); err != nil {
				return err
			}
			break
		}
		if err != nil {
			return err
		}
		_, err = writer.Write(chunk.Data)
	}
	wg.Wait()
	if err != nil {
		return err
	}
	return stream.SendAndClose(&volume_pb.ChunkID{Id: id})
}

func (s *Server) GetChunk(in *volume_pb.ChunkID, stream volume_pb.VolumeService_GetChunkServer) error {
	reader, writer := io.Pipe()
	var err error
	go func() {
		err = s.getChunk(in.Id, writer)
	}()
	buf := mem.Allocate(64 * 1024)
	defer mem.Free(buf)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		err = stream.Send(&volume_pb.ChunkData{Data: buf[:n]})
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) DeleteChunk(ctx context.Context, in *volume_pb.ChunkID) (*emptypb.Empty, error) {
	err := s.deleteChunk(in.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) CreateVolume() (*Volume, error) {
	var id uint64
	var magic []byte
	for i := 0; i < 3; i++ {
		magic = randutil.RandByte(magicLen)
		resp, err := s.client.RegisterVolume(context.Background(), &master_pb.RegisterVolumeRequest{
			Endpoint: s.localEndpoint,
			Id:       uint64(master_pb.ID_CreateId),
			Magic:    magic,
		})
		if err != nil {
			st, ok := status.FromError(err)
			if !ok || st.Code() != codes.Code(ecode.ErrVolumeCreateConflict) {
				return nil, err
			}
		} else {
			id = resp.Id
			break
		}
	}
	v, err := NewVolume(s.dir, id, s.creator)
	if err != nil {
		return nil, err
	}
	err = v.WriteMagic(magic)
	if err != nil {
		return nil, err
	}
	_, err = s.volumeCount.WriteMax(id)
	if err != nil {
		return nil, err
	}
	s.volumeMap.Store(id, v)
	return v, nil
}

func (s *Server) getWritableVolume() *Volume {
	v, _ := s.queue.Get()
	return v.(*Volume)
}

func (s *Server) getVolume(id uint64) (*Volume, error) {
	v, ok := s.volumeMap.Load(id)
	if !ok {
		return nil, errors.ErrVolumeNotFound
	}
	return v.(*Volume), nil
}

func (s *Server) putChunk(reader io.ReadCloser) (string, error) {
	defer func() {
		_ = reader.Close()
	}()
	v := s.getWritableVolume()
	id, err := v.Write(reader)
	if err != nil {
		return "", err
	}
	if v.Size() >= maxVolumeSize {
		err := v.Close()
		if err != nil {
			s.logger.Errorf("close volume failed: %v", err)
		}
		v, err = s.CreateVolume()
		if err != nil {
			s.logger.Errorf("create volume failed: %v", err)
		} else {
			_ = s.queue.Put(v)
		}
	} else {
		_ = s.queue.Put(v)
	}
	return id, nil
}

func (s *Server) getChunk(id string, writer io.WriteCloser) error {
	defer func() {
		_ = writer.Close()
	}()
	volumeID, needleID, err := splitVolumeID(id)
	if err != nil {
		return err
	}
	v, err := s.getVolume(volumeID)
	if err != nil {
		return err
	}
	return v.Read(needleID, writer)
}

func (s *Server) deleteChunk(id string) error {
	volumeID, needleID, err := splitVolumeID(id)
	if err != nil {
		return err
	}
	v, err := s.getVolume(volumeID)
	if err != nil {
		return err
	}
	return v.Delete(needleID)
}

func splitVolumeID(id string) (uint64, needle.Id, error) {
	part := strings.SplitN(id, "-", 2)
	if len(part) != 2 {
		return 0, 0, errors.ErrVolumeIDInvalid
	}
	volumeID, err := strconv.ParseUint(part[0], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	nl, err := strconv.ParseUint(part[1], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return volumeID, needle.Id(nl), nil
}
