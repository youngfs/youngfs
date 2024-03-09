package uuid

import (
	"context"
	"github.com/google/uuid"
)

type SimpleUUID struct {
}

func NewSimpleUUID() *SimpleUUID {
	return &SimpleUUID{}
}

func (u *SimpleUUID) String(ctx context.Context) (string, error) {
	return uuid.NewString(), nil
}

func (u *SimpleUUID) Bytes(ctx context.Context) ([]byte, error) {
	return []byte(uuid.NewString()), nil
}
