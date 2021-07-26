package customerror

type TooManyArgsError struct {
}

func NewTooManyArgsError() *TooManyArgsError {
	return &TooManyArgsError{}
}

func (ce *TooManyArgsError) Error() string {
	return "Too many Arguments!"
}
