package api

import (
	"github.com/go-playground/validator/v10"
)

const (
	ValidationErrorStatus = 400
	ValidationErrorCode   = "ValidationError"
	ValidationErrorTitle  = "Validation error"
)

type ValidationError struct {
	*BaseError

	Violations []validator.FieldError
}

var _ Error = &ValidationError{}
var _ Errorer = &ValidationError{}

type ValidationViolationError struct {
	*BaseError
	Violation validator.FieldError
}

var _ Error = &ValidationViolationError{}

func (e *ValidationViolationError) Detail() string {
	return e.Violation.Error()
}

func NewValidationError(err validator.ValidationErrors) *ValidationError {
	return &ValidationError{
		BaseError: NewBaseError(
			err,
			ValidationErrorStatus,
			ValidationErrorCode,
			ValidationErrorTitle,
		),
		Violations: err,
	}
}

func (e *ValidationError) APIErrors() []Error {
	apiErrors := []Error{}
	for _, violation := range e.Violations {
		apiErrors = append(apiErrors, &ValidationViolationError{
			BaseError: NewBaseError(
				e.originalError,
				e.status,
				e.code,
				e.title,
			),
			Violation: violation,
		})
	}
	return apiErrors
}
