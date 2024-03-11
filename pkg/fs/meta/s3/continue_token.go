package s3

import (
	"context"
	"fmt"
	"github.com/youngfs/youngfs/pkg/errors"
	"github.com/youngfs/youngfs/pkg/kv"
)

func (s *S3) NewContinueToken(ctx context.Context, prefix, startAfter string) (string, error) {
	continueToken, err := s.idGen.String(ctx)
	if err != nil {
		return "", err
	}
	err = s.continueKv.Put(ctx, []byte(fmt.Sprintf("%s-%s", prefix, continueToken)), []byte(startAfter))
	if err != nil {
		return "", err
	}
	return continueToken, nil
}

func (s *S3) GetContinueToken(ctx context.Context, prefix, continueToken string) (string, error) {
	val, err := s.continueKv.Get(ctx, []byte(fmt.Sprintf("%s-%s", prefix, continueToken)))
	if err != nil {
		if errors.Is(err, kv.ErrKeyNotFound) {
			return "", errors.ListObjectsInvalidContinueToken
		}
		return "", err
	}
	return string(val), nil
}
