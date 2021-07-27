package cmd

import (
	"fmt"
	"time"

	"github.com/Zanos420/bcethabot/src/commands"
	"github.com/bwmarrin/discordgo"
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
	p, err := discordgo.SnowflakeTimestamp(ctx.Message.ID)
	diff := time.Until(p)
	ping := diff.Milliseconds()
	ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf("🏓Pong! (Took %v ms)", -ping))
	return
}
