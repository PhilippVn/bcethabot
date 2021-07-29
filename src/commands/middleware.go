package commands

// Middleware Module
// Functional Interface -> Each mittleware will be excecuted
type Middleware interface {
	Exec(ctx *Context, cmd Command) (next bool, err error)   // if next false, excec stopped else next middleware or back to command handler (actual command)
	ExecDM(ctx *Context, cmd Command) (next bool, err error) // Exec for Dm commands
}
