package Error

/* ApplicationError is an error type that should be logged, not provided to the user. */
type ApplicationError interface {
	/* Error returns the string representation of the error, as mandated by the error interface. */
	Error() string

	/* IsApplicationError is an empty placeholder differentiation method due to Go's implicit interface requirement. */
	isApplicationError()
}

/* applicationError is an error type that should be logged, not provided to the user. */
type applicationError struct {
	err string
}

/* NewApplicationError returns a newly initialised ApplicationError. */
func NewApplicationError(err string) ApplicationError {
	return &applicationError{err: err}
}

/* Error returns the string representation of the error, as mandated by the error interface. */
func (e *applicationError) Error() string {
	return e.err
}

/* IsApplicationError is an empty placeholder differentiation method due to Go's implicit interface requirement. */
func (e *applicationError) isApplicationError() {}
