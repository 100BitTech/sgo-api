package base

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/samber/oops"
)

type HttpError struct {
	code int
	error
}

func NewHttpError(code int, err error) error {
	return oops.Wrap(&HttpError{code, err})
}

func NewHttpErrorf(code int, format string, args ...any) error {
	return oops.Wrap(&HttpError{code, fmt.Errorf(format, args...)})
}

func (e *HttpError) Unwrap() error {
	return e.error
}

func (e *HttpError) Code() int {
	return e.code
}

func NewBadRequestError(err error) error {
	return NewHttpError(http.StatusBadRequest, err)
}

func NewBadRequestErrorf(format string, args ...any) error {
	return NewHttpErrorf(http.StatusBadRequest, format, args...)
}

func NewServiceError(err error) error {
	return NewHttpError(http.StatusInternalServerError, err)
}

func NewServiceErrorf(format string, args ...any) error {
	return NewHttpErrorf(http.StatusInternalServerError, format, args...)
}

func NewServiceUnavailableError(err error) error {
	return NewHttpError(http.StatusServiceUnavailable, err)
}

func NewServiceUnavailableErrorf(format string, args ...any) error {
	return NewHttpErrorf(http.StatusServiceUnavailable, format, args...)
}

type MultiError struct {
	errs []error
}

func NewMultiError(errs ...error) *MultiError {
	return &MultiError{errs}
}

func (e *MultiError) IsZero() bool {
	return len(e.errs) == 0
}

func (e *MultiError) Error() string {
	ss := []string{}
	for _, err := range e.errs {
		ss = append(ss, err.Error())
	}
	return strings.Join(ss, "\n")
}

func (e *MultiError) Erros() []error {
	return e.errs
}
