package domain

import "fmt"

type ErrorKind int

const (
	NotFound ErrorKind = iota
	InternalError
)

type ServiceError struct {
	Kind    ErrorKind
	Message error
}

func (err *ServiceError) Error() string {
	if err.Kind == NotFound {
		return "not found"
	}
	if err.Kind == InternalError {
		return fmt.Sprintf("internal error: %v", err.Message)
	}
	panic("unknown error kind")
}

func NewNotFound() *ServiceError {
	return &ServiceError{
		Kind: NotFound,
	}
}

func NewInternalError(message error) *ServiceError {
	return &ServiceError{
		Kind:    InternalError,
		Message: message,
	}
}
