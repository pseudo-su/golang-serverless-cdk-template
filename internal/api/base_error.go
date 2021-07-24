package api

import (
	"fmt"
)

const (
	InternalServerError    = "InternalServerError"
	RequestValidationError = "RequestValidationError"
)

type Error interface {
	error

	Status() int
	Code() string
	Title() string
	Detail() string

	OriginalError() error
}

type Errorer interface {
	error

	APIErrors() []Error
}

func NewBaseError(err error, status int, code string, title string) *BaseError {
	return &BaseError{
		originalError: err,
		status:        status,
		code:          code,
		title:         title,
	}
}

type BaseError struct {
	originalError error
	status        int
	code          string
	title         string
}

var _ Error = &BaseError{}
var _ Errorer = &BaseError{}

func (e *BaseError) Status() int {
	return e.status
}

func (e *BaseError) Code() string {
	return e.code
}

func (e *BaseError) Title() string {
	return e.title
}

func (e *BaseError) Detail() string {
	return e.title
}

func (e *BaseError) OriginalError() error {
	return e.originalError
}

func (e *BaseError) APIErrors() []Error {
	return []Error{e}
}

func (e *BaseError) Error() string {
	return fmt.Sprintf("[%s %d] %s %s: %x", e.Code(), e.Status(), e.Title(), e.Detail(), e.originalError.Error())
}
