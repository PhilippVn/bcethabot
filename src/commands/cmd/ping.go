package cmd

// Ping Command Module
import (
	"fmt"
	"time"

	"github.com/Zanos420/bcethabot/src/commands"
)

type CmdPing struct{}

func NewCmdPing() *CmdPing {
	return &CmdPing{}
}

func (c *CmdPing) Invokes() []string {
	return []string{"ping", "p"} // Invokes and alias
}
func (c *CmdPing) Description() string {
	return "Pong!"
}
func (c *CmdPing) PermissionsNeeded() bool {
	return false
}

func (c *CmdPing) CooldownLocked() bool {
	return true
}

func (c *CmdPing) Exec(ctx *commands.Context) (err error) {
	start := time.Now()
	msg, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "🏓Pong!")
	ping := time.Since(start)
	if err != nil {
		return
	}
	_, err = ctx.Session.ChannelMessageEdit(msg.ChannelID, msg.ID, fmt.Sprintf("🏓Pong! (Took %v ms)", ping.Milliseconds()))
	return
}

func (c *CmdPing) ExecDM(ctx *commands.Context) (err error) {
	start := time.Now()
	msg, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "🏓Pong!")
	ping := time.Since(start)
	if err != nil {
		return
	}
	_, err = ctx.Session.ChannelMessageEdit(msg.ChannelID, msg.ID, fmt.Sprintf("🏓Pong! (Took %v ms)", ping.Milliseconds()))
	return
}
