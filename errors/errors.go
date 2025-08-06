// Package errors / implement public module for service errors
package errors

import (
	"context"
	"errors"
)

type ServiceErrorKind string

const (
	UserVisible ServiceErrorKind = "user_visible"
)

type ServiceError struct {
	Kind ServiceErrorKind `json:"kind"`
	Msg  string           `json:"error"`
}

func NewServiceError(ctx context.Context, msg string, kind ServiceErrorKind) error {
	return &ServiceError{
		Kind: kind,
		Msg:  msg,
	}
}

func Unwrap(err error) *ServiceError {
	var sError *ServiceError

	if errors.As(err, &sError) {
		return sError
	}

	return nil
}

func (e *ServiceError) Error() string {
	if e == nil {
		return ""
	}

	return e.Msg
}
