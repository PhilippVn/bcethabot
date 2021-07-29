package customerror

type InvalidConfigurationError struct{}

func NewInvalidConfigurationError() *InvalidConfigurationError {
	return &InvalidConfigurationError{}
}

func (ce *InvalidConfigurationError) Error() string {
	return "Bot hasnt been configurated correctly!"
}
