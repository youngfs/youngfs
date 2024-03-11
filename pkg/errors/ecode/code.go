package ecode

type Code int

// 100 error code
const (
	ErrKvNotFound Code = 1000 + iota
	ErrChunkNotFound
	ErrNeedleNotFound
)

// 200  error code
const (
	ErrNone Code = 2000 + iota
	ErrCreated
)

// 400  error code
const (
	ErrInvalidPath Code = 4000 + iota
	ErrInvalidDelete
	ErrIllegalObjectName
	ErrIllegalBucketName
	ErrRouter
	ErrChunkNotExist
	ErrContentEncoding
	ErrObjectNotFound
	ErrInvalidContinueToken
)

// 500  error code
const (
	ErrKvSever Code = 5000 + iota
	ErrNonApiError
	ErrProto
	ErrFSServer
	ErrEngineMaster
	ErrEngineChunk
	ErrChunkMisalignment
	ErrMaster
	ErrVolumeMagic
	ErrVolumeCreateConflict
	ErrVolumeNotFound
	ErrVolumeWrite
	ErrVolumeRead
	ErrInvalidNeedle
	ErrVolumeIDInvalid
)
