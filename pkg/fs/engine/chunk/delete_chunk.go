package chunk

import (
	"context"
	"github.com/youngfs/youngfs/pkg/chunk/pb/volume_pb"
	"github.com/youngfs/youngfs/pkg/errors"
)

func (e *Engine) DeleteChunk(ctx context.Context, id string) error {
	volumeId, _, err := ParseNeedle(id)
	if err != nil {
		return err
	}
	client, err := e.getVolumeClientFromId(volumeId)
	if err != nil {
		return errors.ErrEngineChunk.WarpErr(err)
	}
	_, err = client.DeleteChunk(ctx, &volume_pb.ChunkID{Id: id})
	if err != nil {
		return errors.ErrEngineChunk.WarpErr(err)
	}
	return nil
}
