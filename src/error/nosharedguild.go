package customerror

type NoSharedGuildError struct{}

func NewNoSharedGuildError() *NoSharedGuildError {
	return &NoSharedGuildError{}
}

func (ce *NoSharedGuildError) Error() string {
	return "No shared Guild!"
}
