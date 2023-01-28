package errors

type ErrorCode int

// 100 error code
const (
	errKvNotFound ErrorCode = 1000 + iota
)

// 200  error code
const (
	errNone ErrorCode = 2000 + iota
	errCreated
)

// 400  error code
const (
	errInvalidPath ErrorCode = 4000 + iota
	errInvalidDelete
	errIllegalObjectName
	errIllegalSetName
	errIllegalSetRules
	errIllegalObjectSize
	errRouter
	errRecoverFailed
	errObjectNotExist
	errSetRulesNotExist
)

// 500  error code
const (
	errKvSever ErrorCode = 5000 + iota
	errNonApiError
	errProto
	errSeaweedFSMaster
	errSeaweedFSVolume
	errRedisSync
	errServer
)
