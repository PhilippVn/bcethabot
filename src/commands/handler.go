package commands

// Command Handler Module that manages Commands and Middlewares through a message event
import (
	"fmt"
	"strings"

	customerror "github.com/Zanos420/bcethabot/src/error"
	"github.com/bwmarrin/discordgo"
)

// main class that will be used for parsing and excecuting commands
type CommandHandler struct {
	Prefix string

	CmdInstances []Command          // Using Command Interface so all Commands can be saved
	cmdMap       map[string]Command // mapping invoke/alias to Command
	middlewares  []Middleware

	OnError func(err error, ctx *Context) // Public Error fun: Will be Excecuted when a Error occurrs
}

func NewCommandHandler(Prefix string) *CommandHandler {
	return &CommandHandler{
		Prefix:       Prefix,
		CmdInstances: make([]Command, 0),
		cmdMap:       make(map[string]Command),
		middlewares:  make([]Middleware, 0),
		OnError: func(err error, ctx *Context) { // default error funtion
			ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf("Command Excecution failed: %s", err.Error()))
		},
	}
}

func (c *CommandHandler) RegisterCommand(cmd Command) {
	c.CmdInstances = append(c.CmdInstances, cmd) // create new array -> fixed size
	for _, invoke := range cmd.Invokes() {       // register all invokes (alias too)
		c.cmdMap[invoke] = cmd
	}
}

func (c *CommandHandler) RegisterMiddleware(mw Middleware) {
	c.middlewares = append(c.middlewares, mw) // fifo principle: first middleware added is the first being excecuted
}

func (c *CommandHandler) HandleMessage(s *discordgo.Session, e *discordgo.MessageCreate) {
	// Check if the author is a Bot or if the msg statrts
	if e.Author.ID == s.State.User.ID || e.Author.Bot || !strings.HasPrefix(e.Content, c.Prefix) {
		return
	}
	// split Message by Prefix and whitespaces -> invoke and arg are seperated
	split := strings.Split(e.Content[len(c.Prefix):], " ")

	if len(split) < 1 {
		return
	}
	// split Command in invoke and args
	invoke := strings.ToLower(split[0])
	args := split[1:]
	// Check if valid Command
	cmd, ok := c.cmdMap[invoke]
	if !ok || cmd == nil {
		_, err := s.ChannelMessageSend(e.GuildID, customerror.NewCommandNotFoundError().Error())
		if err != nil {
			return
		}
		return
	}

	ctx := &Context{
		Session: s,
		Args:    args,
		Message: e.Message,
		Handler: c,
	}

	if ctx.Message.GuildID == "" {
		guild := ctx.Session.State.Guilds[0]
		for _, member := range guild.Members {
			fmt.Println(member.User.String())
		}
		//Dm command
		// Excecute Middlewares before command
		for _, mw := range c.middlewares {
			next, err := mw.ExecDM(ctx, cmd)
			if err != nil {
				c.OnError(err, ctx) // Call Command Error Handler
				return              // Stop Command Excecution
			}
			if !next {
				return // no further processing
			}
		}

		// Excecute the Command itsself
		if err := cmd.ExecDM(ctx); err != nil {
			c.OnError(err, ctx)
		}
	} else {
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
}
