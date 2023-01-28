package errors

import (
	"github.com/pkg/errors"
)

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func WithMessage(err error, message string) error {
	return errors.WithMessage(err, message)
}

func WithStack(err error) error {
	return errors.WithStack(err)
}

func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

func Wrapf(err error, format string, args ...any) error {
	return errors.Wrapf(err, format, args...)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func New(message string) error {
	return errors.New(message)
}

func Errorf(format string, args ...any) error {
	return errors.Errorf(format, args...)
}

func Cause(err error) error {
	return errors.Cause(err)
}
