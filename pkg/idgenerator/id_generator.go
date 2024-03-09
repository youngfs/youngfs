package idgenerator

import "context"

type IDGenerator interface {
	String(ctx context.Context) (string, error)
	Bytes(ctx context.Context) ([]byte, error)
}
