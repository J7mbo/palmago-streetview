package Error

/* UserError is an error that should be displayed to the user once propagated back to the controller. */
type UserError struct {
	/* Code is a UNIQUE error code key that should be mapped in the GrpcErrorMapper. */
	Code string
	/* Err is the error string to be displayed to the user. */
	Err string
}

/* Error returns the string representation of the error, as mandated by the error interface. */
func (e UserError) Error() string {
	return e.Err
}
