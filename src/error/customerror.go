package customerror

type CustomError struct {
	errormsg string
}

func NewCustomError(errormessage string) *CustomError {
	return &CustomError{errormsg: errormessage}
}

func (ce *CustomError) Error() string {
	return ce.errormsg
}
