package apperrors

import (
	"errors"
	"fmt"
)

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

var (
	ErrNotFound         = errors.New("resource not found")
	ErrPermissionDenied = errors.New("permission denied")
	ErrValidation       = errors.New("validation error")
	ErrInternal         = errors.New("internal server error")
	ErrConflict         = errors.New("resource conflict")
)

func NewNotFound(resource string, id string) error {
	return &AppError{
		Code:    404,
		Message: fmt.Sprintf("%s with ID %s not found", resource, id),
		Err:     ErrNotFound,
	}
}

func NewPermissionDenied(msg string) error {
	return &AppError{
		Code:    403,
		Message: msg,
		Err:     ErrPermissionDenied,
	}
}

func NewValidationError(msg string) error {
	return &AppError{
		Code:    400,
		Message: msg,
		Err:     ErrValidation,
	}
}
