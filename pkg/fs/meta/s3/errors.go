package s3

import "errors"

var (
	ErrS3ObjectNotFound                  = errors.New("s3 object not found")
	ErrS3ListObjectsInvalidContinueToken = errors.New("s3 list objects invalid continue token")
)
