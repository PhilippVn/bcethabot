package customerror

type NotEnoughArgsError struct {
}

func NewNotEnoughArgsError() *NotEnoughArgsError {
	return &NotEnoughArgsError{}
}

func (ce *NotEnoughArgsError) Error() string {
	return "Missing Arguments!"
}
