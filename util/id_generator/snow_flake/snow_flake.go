package snow_flake

import (
	"github.com/bwmarrin/snowflake"
	"youngfs/log"
)

type SnowFlake struct {
	node *snowflake.Node
}

func NewSnowFlake(id int64) *SnowFlake {
	node, err := snowflake.NewNode(id)
	if err != nil {
		log.Fatalf("new snow flake node err: %+v", err)
	}
	return &SnowFlake{
		node: node,
	}
}

func (g *SnowFlake) Generate() (string, error) {
	return g.node.Generate().String(), nil
}
