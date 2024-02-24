package domain

import (
	"strings"
)

// NotFoundErr is an error type for when a resource is not found.
type NotFoundErr struct {
	Msg string
}

func (e *NotFoundErr) Error() string {
	return e.Msg
}

// ValidationErr is an error type when a validation fails.
type ValidationErr struct {
	Err error
}

func (e *ValidationErr) Error() string {
	return e.Err.Error()
}

// CompositeErr is an aggregate error type for multiple errors.
type CompositeErr struct {
	Errs []error
}

func (e *CompositeErr) AddError(err error) {
	e.Errs = append(e.Errs, err)
}

func (e *CompositeErr) Error() string {
	errStrs := make([]string, len(e.Errs))
	for idx, err := range e.Errs {
		errStrs[idx] = err.Error()
	}
	return strings.Join(errStrs, "; ")
}

func (e *CompositeErr) IsEmpty() bool {
	return len(e.Errs) == 0
}
