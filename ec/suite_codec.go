package ec

import (
	"context"
	"github.com/golang/protobuf/proto"
	"icesos/command/vars"
	"icesos/ec/ec_pb"
	"icesos/errors"
	"icesos/full_path"
	"icesos/log"
	"icesos/set"
)

func (suite *Suite) toPb() *ec_pb.Suite {
	if suite == nil {
		return nil
	}

	ShardPb := make([]*ec_pb.Shard, len(suite.Shards))
	for i, u := range suite.Shards {
		ShardPb[i] = u.toPb()
	}

	return &ec_pb.Suite{
		ECid:   suite.ECid,
		Fid:    suite.Fid,
		Shards: ShardPb,
	}
}

func (shard *Shard) toPb() *ec_pb.Shard {
	if shard == nil {
		return nil
	}

	fragsPb := make([]*ec_pb.Frag, len(shard.Frags))
	for i, u := range shard.Frags {
		fragsPb[i] = u.toPb()
	}

	return &ec_pb.Shard{
		Host:  shard.Host,
		Frags: fragsPb,
		Md5:   shard.Md5,
	}
}

func (frag *Frag) toPb() *ec_pb.Frag {
	if frag == nil {
		return nil
	}

	return &ec_pb.Frag{
		FullPath: string(frag.FullPath),
		Set:      string(frag.Set),
		FileSize: frag.FileSize,
		Fid:      frag.Fid,
	}
}

func suitePbToInstance(pb *ec_pb.Suite) *Suite {
	if pb == nil {
		return nil
	}

	shards := make([]Shard, len(pb.Shards))
	for i, u := range pb.Shards {
		if u == nil {
			continue
		}
		shards[i] = *shardPbToInstance(u)
	}

	if pb.Shards == nil {
		shards = nil
	}

	return &Suite{
		ECid:   pb.ECid,
		Fid:    pb.Fid,
		Shards: shards,
	}
}

func shardPbToInstance(pb *ec_pb.Shard) *Shard {
	if pb == nil {
		return nil
	}

	frags := make([]Frag, len(pb.Frags))
	for i, u := range pb.Frags {
		if u == nil {
			continue
		}
		frags[i] = *fragPbToInstance(u)
	}

	if pb.Frags == nil {
		frags = nil
	}

	return &Shard{
		Host:  pb.Host,
		Frags: frags,
		Md5:   pb.Md5,
	}
}

func fragPbToInstance(pb *ec_pb.Frag) *Frag {
	if pb == nil {
		return nil
	}

	return &Frag{
		FullPath: full_path.FullPath(pb.FullPath),
		Set:      set.Set(pb.Set),
		Fid:      pb.Fid,
		FileSize: pb.FileSize,
	}
}

func (suite *Suite) EncodeProto(ctx context.Context) ([]byte, error) {
	message := suite.toPb()
	b, err := proto.Marshal(message)
	if err != nil {
		log.Errorw("encode suite proto error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error())
		err = errors.GetAPIErr(errors.ErrProto)
	}
	return b, err
}

func (shard *Shard) EncodeProto(ctx context.Context) ([]byte, error) {
	message := shard.toPb()
	b, err := proto.Marshal(message)
	if err != nil {
		log.Errorw("encode shard proto error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error())
		err = errors.GetAPIErr(errors.ErrProto)
	}
	return b, err
}

func (frag *Frag) EncodeProto(ctx context.Context) ([]byte, error) {
	message := frag.toPb()
	b, err := proto.Marshal(message)
	if err != nil {
		log.Errorw("encode frag proto error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error())
		err = errors.GetAPIErr(errors.ErrProto)
	}
	return b, err
}

func DecodeSuiteProto(ctx context.Context, b []byte) (*Suite, error) {
	message := &ec_pb.Suite{}
	if err := proto.Unmarshal(b, message); err != nil {
		log.Errorw("decode suite proto error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error())
		return nil, errors.GetAPIErr(errors.ErrProto)
	}
	return suitePbToInstance(message), nil
}

func DecodeShardProto(ctx context.Context, b []byte) (*Shard, error) {
	message := &ec_pb.Shard{}
	if err := proto.Unmarshal(b, message); err != nil {
		log.Errorw("decode shard proto error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error())
		return nil, errors.GetAPIErr(errors.ErrProto)
	}
	return shardPbToInstance(message), nil
}

func DecodeFragProto(ctx context.Context, b []byte) (*Frag, error) {
	message := &ec_pb.Frag{}
	if err := proto.Unmarshal(b, message); err != nil {
		log.Errorw("decode frag proto error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error())
		return nil, errors.GetAPIErr(errors.ErrProto)
	}
	return fragPbToInstance(message), nil
}
