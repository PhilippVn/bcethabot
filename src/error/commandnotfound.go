package customerror

type CommandNotFoundError struct {
}

func NewCommandNotFoundError() *CommandNotFoundError {
	return &CommandNotFoundError{}
}

func (ce *CommandNotFoundError) Error() string {
	return "Command not found!"
}
