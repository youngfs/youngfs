package chunk

import (
	"context"
	"github.com/youngfs/youngfs/pkg/chunk/pb/volume_pb"
	"github.com/youngfs/youngfs/pkg/errors"
	"github.com/youngfs/youngfs/pkg/util/mem"
	"io"
)

func (e *Engine) PutChunk(ctx context.Context, reader io.Reader, endpoints ...string) (string, error) {
	client, err := e.getVolumeClientFromEndpoints(endpoints...)
	if err != nil {
		return "", err
	}
	stream, err := client.PutChunk(ctx)
	if err != nil {
		return "", errors.ErrEngineChunk.WarpErr(err)
	}
	buf := mem.Allocate(128 * 1024)
	defer mem.Free(buf)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", errors.ErrEngineChunk.WarpErr(err)
		}
		err = stream.Send(&volume_pb.ChunkData{Data: buf[:n]})
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return "", errors.ErrEngineChunk.WarpErr(err)
	}
	return resp.Id, nil
}
