package errors

import (
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"io"
	"net/http"
)

type APIError struct {
	ErrorCode
	HTTPStatusCode int
	Description    string
	*stack
}

func (e *APIError) Error() string {
	return fmt.Sprintf("HTTP status Code: %d. Error code: %d. Description: %s.", e.HTTPStatusCode, e.ErrorCode, e.Description)
}

func (e *APIError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = io.WriteString(s, e.Error())
			if e.stack != nil {
				e.stack.Format(s, verb)
			}
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
	ErrKvNotFound  = &APIError{ErrorCode: errKvNotFound, HTTPStatusCode: http.StatusContinue, Description: "Kv not found"}
	ErrEcNotFinish = &APIError{ErrorCode: errECNotFinish, HTTPStatusCode: http.StatusContinue, Description: "Ec not finish"}
	// 200
	ErrNone    = &APIError{ErrorCode: errNone, HTTPStatusCode: http.StatusOK, Description: "Request succeeded"}
	ErrCreated = &APIError{ErrorCode: errCreated, HTTPStatusCode: http.StatusCreated, Description: "Created succeeded"}
	// 400
	ErrInvalidPath       = &APIError{ErrorCode: errInvalidPath, HTTPStatusCode: http.StatusNotFound, Description: "The file path is not valid"}
	ErrInvalidDelete     = &APIError{ErrorCode: errInvalidDelete, HTTPStatusCode: http.StatusBadRequest, Description: "There are files in the folder and cannot be deleted recursively"}
	ErrIllegalObjectName = &APIError{ErrorCode: errIllegalObjectName, HTTPStatusCode: http.StatusBadRequest, Description: "Illegal object name"}
	ErrIllegalSetName    = &APIError{ErrorCode: errIllegalSetName, HTTPStatusCode: http.StatusBadRequest, Description: "Illegal set name"}
	ErrIllegalSetRules   = &APIError{ErrorCode: errIllegalSetRules, HTTPStatusCode: http.StatusBadRequest, Description: "Illegal set rules"}
	ErrIllegalObjectSize = &APIError{ErrorCode: errIllegalObjectSize, HTTPStatusCode: http.StatusBadRequest, Description: "Illegal object size"}
	ErrRouter            = &APIError{ErrorCode: errRouter, HTTPStatusCode: http.StatusBadRequest, Description: "Router problem"}
	ErrRecoverFailed     = &APIError{ErrorCode: errRecoverFailed, HTTPStatusCode: http.StatusBadRequest, Description: "Object recover failed"}
	ErrObjectNotExist    = &APIError{ErrorCode: errObjectNotExist, HTTPStatusCode: http.StatusNotFound, Description: "Object not exist"}
	ErrRulesNotExist     = &APIError{ErrorCode: errRulesNotExist, HTTPStatusCode: http.StatusNotFound, Description: "set rules not exist"}
	// 500
	ErrKvSever         = &APIError{ErrorCode: errKvSever, HTTPStatusCode: http.StatusInternalServerError, Description: "Key-value database error"}
	ErrNonApiErr       = &APIError{ErrorCode: errNonApiError, HTTPStatusCode: http.StatusInternalServerError, Description: "Non api error return"}
	ErrProto           = &APIError{ErrorCode: errProto, HTTPStatusCode: http.StatusInternalServerError, Description: "ProtoBuf error"}
	ErrSeaweedFSMaster = &APIError{ErrorCode: errSeaweedFSMaster, HTTPStatusCode: http.StatusInternalServerError, Description: "SeaweedFS master server error"}
	ErrSeaweedFSVolume = &APIError{ErrorCode: errSeaweedFSVolume, HTTPStatusCode: http.StatusInternalServerError, Description: "SeaweedFS volume server error"}
	ErrRedisSync       = &APIError{ErrorCode: errRedisSync, HTTPStatusCode: http.StatusInternalServerError, Description: "Redis lock server error"}
	ErrServer          = &APIError{ErrorCode: errServer, HTTPStatusCode: http.StatusInternalServerError, Description: "youngfs server error"}
)

func (e *APIError) Wrap(msg string) error {
	return errors.Wrap(e, msg)
}

func (e *APIError) WithMessage(msg string) error {
	return errors.WithMessage(e, msg)
}

func (e *APIError) Wrapf(format string, args ...any) error {
	return errors.Wrapf(e, format, args...)
}

func (e *APIError) WithStack() *APIError {
	return &APIError{
		ErrorCode:      e.ErrorCode,
		HTTPStatusCode: e.HTTPStatusCode,
		Description:    e.Description,
		stack:          callers(),
	}
}

func (e *APIError) WrapErr(err error) error {
	return multierr.Append(e.WithStack(), err)
}

func (e *APIError) WrapErrNoStack(err error) error {
	return multierr.Append(e, err)
}

func (e *APIError) IsServerErr() bool {
	return (int(e.ErrorCode) / 5000) > 0
}

func IsKvNotFound(err error) bool {
	return Is(err, ErrKvNotFound)
}
