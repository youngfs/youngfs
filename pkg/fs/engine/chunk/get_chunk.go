package chunk

import (
	"context"
	"github.com/youngfs/youngfs/pkg/chunk/pb/volume_pb"
	"github.com/youngfs/youngfs/pkg/chunk/volume"
	"github.com/youngfs/youngfs/pkg/chunk/volume/needle"
	"github.com/youngfs/youngfs/pkg/errors"
	"io"
)

func (e *Engine) GetChunk(ctx context.Context, id string) (io.ReadCloser, error) {
	volumeId, _, err := ParseNeedle(id)
	if err != nil {
		return nil, err
	}
	client, err := e.getVolumeClientFromId(volumeId)
	if err != nil {
		return nil, err
	}
	stream, err := client.GetChunk(ctx, &volume_pb.ChunkID{Id: id})
	if err != nil {
		return nil, errors.ErrEngineChunk.WarpErr(err)
	}
	reader, writer := io.Pipe()
	chunk, err := stream.Recv()
	if err == io.EOF {
		if err := writer.Close(); err != nil {
			return nil, errors.ErrEngineChunk.WarpErr(err)
		}
	} else if err != nil {
		return nil, errors.ErrEngineChunk.WarpErr(err)
	}
	go func() {
		_, err = writer.Write(chunk.Data)
		for {
			chunk, err := stream.Recv()
			if err == io.EOF {
				if err := writer.Close(); err != nil {
					return
				}
				break
			} else if err != nil {
				_ = writer.CloseWithError(err)
				return
			}
			_, err = writer.Write(chunk.Data)
		}
	}()
	return reader, nil
}

func ParseNeedle(id string) (uint64, needle.Id, error) {
	volumeID, nl, err := volume.SplitVolumeID(id)
	if err != nil {
		return 0, 0, errors.ErrEngineMaster.WarpErr(err)
	}
	return volumeID, nl, nil
}
