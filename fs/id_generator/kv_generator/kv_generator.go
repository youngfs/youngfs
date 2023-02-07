package kv_generator

import (
	"context"
	"strconv"
	"youngfs/kv"
)

type KvGenerator struct {
	ctx   context.Context
	key   string
	store kv.KvSetStore
}

func NewKvGenerator(key string, store kv.KvSetStore) *KvGenerator {
	return &KvGenerator{
		ctx:   context.Background(),
		key:   key,
		store: store,
	}
}

func (g *KvGenerator) Generate() (string, error) {
	val, err := g.store.Incr(g.ctx, g.key)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(val, 10), nil
}
