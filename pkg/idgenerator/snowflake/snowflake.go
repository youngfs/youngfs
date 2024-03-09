package snowflake

import (
	"context"
	"github.com/bwmarrin/snowflake"
)

type SnowFlake struct {
	node *snowflake.Node
}

func New(id int64) (*SnowFlake, error) {
	node, err := snowflake.NewNode(id)
	if err != nil {
		return nil, err
	}
	return &SnowFlake{
		node: node,
	}, nil
}

func (s *SnowFlake) String(ctx context.Context) (string, error) {
	return s.node.Generate().String(), nil
}

func (s *SnowFlake) Bytes(ctx context.Context) ([]byte, error) {
	return s.node.Generate().Bytes(), nil
}
