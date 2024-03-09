package errors

type Code int

// 100 error code
const (
	errKvNotFound Code = 1000 + iota
)

// 200  error code
const (
	errNone Code = 2000 + iota
	errCreated
)

// 400  error code
const (
	errInvalidPath Code = 4000 + iota
	errInvalidDelete
	errIllegalObjectName
	errIllegalBucketName
	errRouter
	errChunkNotExist
	errContentEncoding
)

// 500  error code
const (
	errKvSever Code = 5000 + iota
	errNonApiError
	errProto
	errFSServer
	errEngineMaster
	errEngineChunk
	errChunkMisalignment
)
