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
	ErrKvNotFound           ErrorCode = 1001
	ErrNone                 ErrorCode = 2001
	ErrCreated              ErrorCode = 2002
	ErrInvalidPath          ErrorCode = 4001
	ErrInvalidDelete        ErrorCode = 4002
	ErrIllegalObjectName    ErrorCode = 4003
	ErrAdminAuthenticate    ErrorCode = 4004
	ErrUserNotExist         ErrorCode = 4005
	ErrUserAuthenticate     ErrorCode = 4006
	ErrSetReadAuthenticate  ErrorCode = 4007
	ErrSetWriteAuthenticate ErrorCode = 4008
	ErrInvalidUserName      ErrorCode = 4009
	ErrKvSever              ErrorCode = 5001
	ErrProto                ErrorCode = 5002
	ErrSeaweedFSMaster      ErrorCode = 5003
	ErrSeaweedFSVolume      ErrorCode = 5004
	ErrRedisSync            ErrorCode = 5005
	ErrServer               ErrorCode = 5006
)

var ErrorCodeResponse = map[ErrorCode]APIError{
	// 200
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
	// 400
	ErrInvalidPath: {
		ErrorCode:      ErrInvalidPath,
		Description:    "The file path is not valid",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrInvalidDelete: {
		ErrorCode:      ErrInvalidDelete,
		Description:    "There are files in the folder and cannot be deleted recursively",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrIllegalObjectName: {
		ErrorCode:      ErrIllegalObjectName,
		Description:    "Illegal object name",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrAdminAuthenticate: {
		ErrorCode:      ErrAdminAuthenticate,
		Description:    "Administrator authority authentication failed",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrUserNotExist: {
		ErrorCode:      ErrUserNotExist,
		Description:    "User does not exist",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrUserAuthenticate: {
		ErrorCode:      ErrUserAuthenticate,
		Description:    "User authority authentication failed",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrSetReadAuthenticate: {
		ErrorCode:      ErrSetReadAuthenticate,
		Description:    "Set read authority authentication failed",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrSetWriteAuthenticate: {
		ErrorCode:      ErrSetWriteAuthenticate,
		Description:    "Set write authority authentication failed",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidUserName: {
		ErrorCode:      ErrInvalidUserName,
		Description:    "Invalid user name",
		HTTPStatusCode: http.StatusBadRequest,
	},
	// 500
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
	ErrServer: {
		ErrorCode:      ErrServer,
		Description:    "icesos server error",
		HTTPStatusCode: http.StatusInternalServerError,
	},
}
