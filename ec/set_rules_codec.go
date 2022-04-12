package ec

import (
	"github.com/golang/protobuf/proto"
	"icesos/ec/ec_pb"
	"icesos/errors"
	"icesos/set"
)

func (setRules *SetRules) toPb() *ec_pb.SetRules {
	if setRules == nil {
		return nil
	}

	return &ec_pb.SetRules{
		Set:             string(setRules.Set),
		Hosts:           setRules.Hosts,
		DataShards:      setRules.DataShards,
		ParityShards:    setRules.ParityShards,
		MAXBlockSize:    setRules.MAXBlockSize,
		ECMode:          setRules.ECMode,
		ReplicationMode: setRules.ReplicationMode,
	}
}

func setRulesPbToInstance(pb *ec_pb.SetRules) *SetRules {
	if pb == nil {
		return nil
	}

	return &SetRules{
		Set:             set.Set(pb.Set),
		Hosts:           pb.Hosts,
		DataShards:      pb.DataShards,
		ParityShards:    pb.ParityShards,
		MAXBlockSize:    pb.MAXBlockSize,
		ECMode:          pb.ECMode,
		ReplicationMode: pb.ReplicationMode,
	}
}

func (setRules *SetRules) EncodeProto() ([]byte, error) {
	message := setRules.toPb()
	b, err := proto.Marshal(message)
	if err != nil {
		err = errors.ErrorCodeResponse[errors.ErrProto]
	}
	return b, err
}

func DecodeSetRulesProto(b []byte) (*SetRules, error) {
	message := &ec_pb.SetRules{}
	if err := proto.Unmarshal(b, message); err != nil {
		return nil, errors.ErrorCodeResponse[errors.ErrProto]
	}
	return setRulesPbToInstance(message), nil
}
