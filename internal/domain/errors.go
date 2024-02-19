package domain

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
