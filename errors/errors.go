package errors

import (
	"fmt"
	"net/http"
)

type ErrorCode int

type APIError struct {
	ErrorCode
	Description    string
	HTTPStatusCode int
}

func (e APIError) Error() string {
	return fmt.Sprintf("HTTP status Code: %d. Error code: %d. Description: %s.", e.HTTPStatusCode, e.ErrorCode, e.Description)
}

const (
	ErrNone            ErrorCode = 2001
	ErrCreated         ErrorCode = 2002
	ErrInvalidPath     ErrorCode = 4001
	ErrInvalidDelete   ErrorCode = 4002
	ErrKvSever         ErrorCode = 5001
	ErrProto           ErrorCode = 5002
	ErrSeaweedFSMaster ErrorCode = 5003
	ErrSeaweedFSVolume ErrorCode = 5004
	ErrRedisSync       ErrorCode = 5005
)

var ErrorCodeResponse = map[ErrorCode]APIError{
	ErrNone: {
		ErrorCode:      ErrNone,
		Description:    "Request succeeded",
		HTTPStatusCode: http.StatusOK,
	},
	ErrCreated: {
		ErrorCode:      ErrCreated,
		Description:    "Created succeeded",
		HTTPStatusCode: http.StatusCreated,
	},
	ErrInvalidPath: {
		ErrorCode:      ErrInvalidPath,
		Description:    "The file path is not valid",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidDelete: {
		ErrorCode:      ErrInvalidDelete,
		Description:    "There are files in the folder and cannot be deleted recursively",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrKvSever: {
		ErrorCode:      ErrKvSever,
		Description:    "Key-value database error",
		HTTPStatusCode: http.StatusInternalServerError,
	},
	ErrProto: {
		ErrorCode:      ErrProto,
		Description:    "ProtoBuf error",
		HTTPStatusCode: http.StatusInternalServerError,
	},
	ErrSeaweedFSMaster: {
		ErrorCode:      ErrSeaweedFSMaster,
		Description:    "SeaweedFS master server error",
		HTTPStatusCode: http.StatusInternalServerError,
	},
	ErrSeaweedFSVolume: {
		ErrorCode:      ErrSeaweedFSVolume,
		Description:    "SeaweedFS volume server error",
		HTTPStatusCode: http.StatusInternalServerError,
	},
	ErrRedisSync: {
		ErrorCode:      ErrRedisSync,
		Description:    "Redis lock server error",
		HTTPStatusCode: http.StatusInternalServerError,
	},
}
