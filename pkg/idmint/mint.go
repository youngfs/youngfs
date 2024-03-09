package idmint

import "context"

type Mint interface {
	String(ctx context.Context) (string, error)
	Bytes(ctx context.Context) ([]byte, error)
}
