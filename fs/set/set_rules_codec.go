package set

import (
	"context"
	"github.com/golang/protobuf/proto"
	"youngfs/errors"
	"youngfs/fs/set/set_pb"
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
		err = errors.ErrProto.WithStack()
	}
	return b, err
}

func DecodeSetRulesProto(ctx context.Context, b []byte) (*SetRules, error) {
	message := &set_pb.SetRules{}
	if err := proto.Unmarshal(b, message); err != nil {
		return nil, errors.ErrProto.WithStack()
	}
	return setRulesPbToInstance(message), nil
}
