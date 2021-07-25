package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// main class that will be used for parsing and excecuting commands
type CommandHandler struct {
	prefix string

	cmdInstances []Command          // Using Command Interface so all Commands can be saved
	cmdMap       map[string]Command // mapping invoke/alias to Command
	middlewares  []Middleware

	OnError func(err error, ctx *Context) // Public Error fun: Will be Excecuted when a Error occurrs
}

func NewCommandHandler(prefix string) *CommandHandler {
	return &CommandHandler{
		prefix:       prefix,
		cmdInstances: make([]Command, 0),
		cmdMap:       make(map[string]Command),
		middlewares:  make([]Middleware, 0),
		OnError:      func(error, *Context) {}, // fallback if nil
	}
}

func (c *CommandHandler) RegisterCommand(cmd Command) {
	c.cmdInstances = append(c.cmdInstances, cmd) // create new array -> fixed size
	for _, invoke := range cmd.Invokes() {       // register all invokes (alias too)
		c.cmdMap[invoke] = cmd
	}
}

func (c *CommandHandler) RegisterMiddleware(mw Middleware) {
	c.middlewares = append(c.middlewares, mw) // fifo principle: first middleware added is the first being excecuted
}

func (c *CommandHandler) HandleMessage(s *discordgo.Session, e *discordgo.MessageCreate) {
	// Check if the author is a Bot or if the msg statrts
	if e.Author.ID == s.State.User.ID || e.Author.Bot || !strings.HasPrefix(e.Content, c.prefix) {
		return
	}
	// split Message by prefix and whitespaces -> invoke and arg are seperated
	split := strings.Split(e.Content[len(c.prefix):], " ")

	if len(split) < 1 {
		return
	}
	// split Command in invoke and args
	invoke := strings.ToLower(split[0])
	args := split[1:]
	// Check if valid Command
	cmd, ok := c.cmdMap[invoke]
	if !ok || cmd == nil {
		return
	}

	ctx := &Context{
		Session: s,
		Args:    args,
		Message: e.Message,
		Handler: c,
	}

	// Excecute Middlewares before command
	for _, mw := range c.middlewares {
		next, err := mw.Exec(ctx, cmd)
		if err != nil {
			c.OnError(err, ctx) // Call Command Error Handler
			return              // Stop Command Excecution
		}
		if !next {
			return // no further processing
		}
	}

	// Excecute the Command itsself
	if err := cmd.Exec(ctx); err != nil {
		c.OnError(err, ctx)
	}
}