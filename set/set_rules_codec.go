package set

import (
	"context"
	"github.com/golang/protobuf/proto"
	"icesos/command/vars"
	"icesos/errors"
	"icesos/log"
	"icesos/set/set_pb"
)

func (setRules *SetRules) toPb() *set_pb.SetRules {
	if setRules == nil {
		return nil
	}

	return &set_pb.SetRules{
		Set:             string(setRules.Set),
		Hosts:           setRules.Hosts,
		DataShards:      setRules.DataShards,
		ParityShards:    setRules.ParityShards,
		MAXShardSize:    setRules.MAXShardSize,
		ECMode:          setRules.ECMode,
		ReplicationMode: setRules.ReplicationMode,
	}
}

func setRulesPbToInstance(pb *set_pb.SetRules) *SetRules {
	if pb == nil {
		return nil
	}

	return &SetRules{
		Set:             Set(pb.Set),
		Hosts:           pb.Hosts,
		DataShards:      pb.DataShards,
		ParityShards:    pb.ParityShards,
		MAXShardSize:    pb.MAXShardSize,
		ECMode:          pb.ECMode,
		ReplicationMode: pb.ReplicationMode,
	}
}

func (setRules *SetRules) EncodeProto(ctx context.Context) ([]byte, error) {
	message := setRules.toPb()
	b, err := proto.Marshal(message)
	if err != nil {
		log.Errorw("encode set rules proto error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error())
		err = errors.GetAPIErr(errors.ErrProto)
	}
	return b, err
}

func DecodeSetRulesProto(ctx context.Context, b []byte) (*SetRules, error) {
	message := &set_pb.SetRules{}
	if err := proto.Unmarshal(b, message); err != nil {
		log.Errorw("decode set rules proto error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error())
		return nil, errors.GetAPIErr(errors.ErrProto)
	}
	return setRulesPbToInstance(message), nil
}
