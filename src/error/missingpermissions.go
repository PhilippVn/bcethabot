package customerror

type MissingPermissionsError struct{}

func NewMissingPermissionsError() *MissingPermissionsError {
	return &MissingPermissionsError{}
}

func (ce *MissingPermissionsError) Error() string {
	return ":no_entry: Missing Permissions!"
}
