package errors

import (
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"io"
	"net/http"
)

type Error struct {
	Code
	HTTPStatusCode int
	Description    string
}

func (e *Error) Error() string {
	return fmt.Sprintf("HTTP status Code: %d. Error code: %d. Description: %s.", e.HTTPStatusCode, e.Code, e.Description)
}

func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = io.WriteString(s, e.Error())
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, e.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", e.Error())
	}
}

var (
	// 100
	// ErrKVNotFound
	// 200
	ErrNone    = &Error{Code: errNone, HTTPStatusCode: http.StatusOK, Description: "Request succeeded"}
	ErrCreated = &Error{Code: errCreated, HTTPStatusCode: http.StatusCreated, Description: "Created succeeded"}
	// 400
	ErrInvalidPath       = &Error{Code: errInvalidPath, HTTPStatusCode: http.StatusNotFound, Description: "The file path is not valid"}
	ErrInvalidDelete     = &Error{Code: errInvalidDelete, HTTPStatusCode: http.StatusBadRequest, Description: "There are files in the folder and cannot be deleted recursively"}
	ErrIllegalObjectName = &Error{Code: errIllegalObjectName, HTTPStatusCode: http.StatusBadRequest, Description: "Illegal object name"}
	ErrIllegalBucketName = &Error{Code: errIllegalBucketName, HTTPStatusCode: http.StatusBadRequest, Description: "Illegal bucket name"}
	ErrRouter            = &Error{Code: errRouter, HTTPStatusCode: http.StatusBadRequest, Description: "Router problem"}
	ErrChunkNotExist     = &Error{Code: errChunkNotExist, HTTPStatusCode: http.StatusNotFound, Description: "Chunk not exist"}
	ErrContentEncoding   = &Error{Code: errContentEncoding, HTTPStatusCode: http.StatusBadRequest, Description: "Content Encoding read error"}
	// 500
	ErrKvSever           = &Error{Code: errKvSever, HTTPStatusCode: http.StatusInternalServerError, Description: "Key-value database error"}
	ErrNonApiErr         = &Error{Code: errNonApiError, HTTPStatusCode: http.StatusInternalServerError, Description: "Non api error return"}
	ErrFSServer          = &Error{Code: errFSServer, HTTPStatusCode: http.StatusInternalServerError, Description: "File system server error"}
	ErrProto             = &Error{Code: errProto, HTTPStatusCode: http.StatusInternalServerError, Description: "ProtoBuf error"}
	ErrEngineMaster      = &Error{Code: errEngineMaster, HTTPStatusCode: http.StatusInternalServerError, Description: "Engine master server error"}
	ErrEngineChunk       = &Error{Code: errEngineChunk, HTTPStatusCode: http.StatusInternalServerError, Description: "Engine chunk server error"}
	ErrChunkMisalignment = &Error{Code: errChunkMisalignment, HTTPStatusCode: http.StatusInternalServerError, Description: "Chunk offset misalignment"}
)

func (e *Error) WithMessage(msg string) error {
	return errors.WithMessage(e, msg)
}

func (e *Error) WithMessagef(format string, args ...any) error {
	return errors.WithMessagef(e, format, args...)
}

func (e *Error) Wrap(msg string) error {
	return errors.Wrap(e, msg)
}

func (e *Error) Wrapf(format string, args ...any) error {
	return errors.Wrapf(e, format, args...)
}

func (e *Error) WarpErr(err error) error {
	return multierr.Append(e, err)
}

func (e *Error) IsServerErr() bool {
	return (int(e.Code) / 5000) > 0
}
