package rules

import (
	"context"
	"github.com/golang/protobuf/proto"
	"youngfs/errors"
	"youngfs/fs/rules/rules_pb"
	"youngfs/fs/set"
)

func (rules *Rules) toPb() *rules_pb.Rules {
	if rules == nil {
		return nil
	}

	return &rules_pb.Rules{
		Set:             string(rules.Set),
		Hosts:           rules.Hosts,
		DataShards:      rules.DataShards,
		ParityShards:    rules.ParityShards,
		MaxShardSize:    rules.MaxShardSize,
		ECMode:          rules.ECMode,
		ReplicationMode: rules.ReplicationMode,
	}
}

func rulesPbToInstance(pb *rules_pb.Rules) *Rules {
	if pb == nil {
		return nil
	}

	return &Rules{
		Set:             set.Set(pb.Set),
		Hosts:           pb.Hosts,
		DataShards:      pb.DataShards,
		ParityShards:    pb.ParityShards,
		MaxShardSize:    pb.MaxShardSize,
		ECMode:          pb.ECMode,
		ReplicationMode: pb.ReplicationMode,
	}
}

func (rules *Rules) EncodeProto(ctx context.Context) ([]byte, error) {
	message := rules.toPb()
	b, err := proto.Marshal(message)
	if err != nil {
		err = errors.ErrProto.WithStack()
	}
	return b, err
}

func DecodeRulesProto(ctx context.Context, b []byte) (*Rules, error) {
	message := &rules_pb.Rules{}
	if err := proto.Unmarshal(b, message); err != nil {
		return nil, errors.ErrProto.WithStack()
	}
	return rulesPbToInstance(message), nil
}
