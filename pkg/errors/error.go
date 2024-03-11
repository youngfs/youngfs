package errors

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/youngfs/youngfs/pkg/errors/ecode"
	"go.uber.org/multierr"
	"io"
	"net/http"
)

type Error struct {
	ecode.Code
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

func (e *Error) ErrorCode() ecode.Code {
	return e.Code
}

var (
	// 100
	// ErrKVNotFound
	// 200
	ErrNone    = &Error{Code: ecode.ErrNone, HTTPStatusCode: http.StatusOK, Description: "Request succeeded"}
	ErrCreated = &Error{Code: ecode.ErrCreated, HTTPStatusCode: http.StatusCreated, Description: "Created succeeded"}
	// 400
	ErrInvalidPath                  = &Error{Code: ecode.ErrInvalidPath, HTTPStatusCode: http.StatusNotFound, Description: "The file path is not valid"}
	ErrInvalidDelete                = &Error{Code: ecode.ErrInvalidDelete, HTTPStatusCode: http.StatusBadRequest, Description: "There are files in the folder and cannot be deleted recursively"}
	ErrIllegalObjectName            = &Error{Code: ecode.ErrIllegalObjectName, HTTPStatusCode: http.StatusBadRequest, Description: "Illegal object name"}
	ErrIllegalBucketName            = &Error{Code: ecode.ErrIllegalBucketName, HTTPStatusCode: http.StatusBadRequest, Description: "Illegal bucket name"}
	ErrRouter                       = &Error{Code: ecode.ErrRouter, HTTPStatusCode: http.StatusBadRequest, Description: "Router problem"}
	ErrChunkNotExist                = &Error{Code: ecode.ErrChunkNotExist, HTTPStatusCode: http.StatusNotFound, Description: "Chunk not exist"}
	ErrContentEncoding              = &Error{Code: ecode.ErrContentEncoding, HTTPStatusCode: http.StatusBadRequest, Description: "Content Encoding read error"}
	ErrObjectNotFound               = &Error{Code: ecode.ErrObjectNotFound, HTTPStatusCode: http.StatusNotFound, Description: "Object not found"}
	ListObjectsInvalidContinueToken = &Error{Code: ecode.ListObjectsInvalidContinueToken, HTTPStatusCode: http.StatusBadRequest, Description: "Invalid continue token"}

	// 500
	ErrKvSever              = &Error{Code: ecode.ErrKvSever, HTTPStatusCode: http.StatusInternalServerError, Description: "Key-value database error"}
	ErrNonApiErr            = &Error{Code: ecode.ErrNonApiError, HTTPStatusCode: http.StatusInternalServerError, Description: "Non api error return"}
	ErrFSServer             = &Error{Code: ecode.ErrFSServer, HTTPStatusCode: http.StatusInternalServerError, Description: "File system server error"}
	ErrProto                = &Error{Code: ecode.ErrProto, HTTPStatusCode: http.StatusInternalServerError, Description: "ProtoBuf error"}
	ErrEngineMaster         = &Error{Code: ecode.ErrEngineMaster, HTTPStatusCode: http.StatusInternalServerError, Description: "Engine master server error"}
	ErrEngineChunk          = &Error{Code: ecode.ErrEngineChunk, HTTPStatusCode: http.StatusInternalServerError, Description: "Engine chunk server error"}
	ErrChunkMisalignment    = &Error{Code: ecode.ErrChunkMisalignment, HTTPStatusCode: http.StatusInternalServerError, Description: "Chunk offset misalignment"}
	ErrMaster               = &Error{Code: ecode.ErrMaster, HTTPStatusCode: http.StatusInternalServerError, Description: "Master server error"}
	ErrVolumeMagic          = &Error{Code: ecode.ErrVolumeMagic, HTTPStatusCode: http.StatusInternalServerError, Description: "Volume check code error"}
	ErrVolumeCreateConflict = &Error{Code: ecode.ErrVolumeCreateConflict, HTTPStatusCode: http.StatusInternalServerError, Description: "Volume create error"}
	ErrVolumeNotFound       = &Error{Code: ecode.ErrVolumeNotFound, HTTPStatusCode: http.StatusNotFound, Description: "Volume not found"}
	ErrVolumeWrite          = &Error{Code: ecode.ErrVolumeWrite, HTTPStatusCode: http.StatusInternalServerError, Description: "Volume write error"}
	ErrVolumeRead           = &Error{Code: ecode.ErrVolumeRead, HTTPStatusCode: http.StatusInternalServerError, Description: "Volume read error"}
	ErrInvalidNeedle        = &Error{Code: ecode.ErrInvalidNeedle, HTTPStatusCode: http.StatusInternalServerError, Description: "Invalid needle"}
	ErrVolumeIDInvalid      = &Error{Code: ecode.ErrVolumeIDInvalid, HTTPStatusCode: http.StatusInternalServerError, Description: "Invalid volume id"}
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
