package commands

// Command Interface
// A command constists of a message split in two parts:
// - prefix
// - Invoke (command name or alias)
type Command interface {
	Invokes() []string       // first element of the slice is the invoke the rest are alias
	Description() string     // short description of the command
	PermissionsNeeded() bool // true if command requires permissions
	Exec(ctx *Context) error // Function that is excecuted when command is issued
}
