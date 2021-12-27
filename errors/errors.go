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
	ErrNone          ErrorCode = 2001
	ErrCreated       ErrorCode = 2002
	ErrInvalidPath   ErrorCode = 4001
	ErrInvalidDelete ErrorCode = 4002
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
}
