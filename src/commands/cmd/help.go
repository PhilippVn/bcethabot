package cmd

import (
	"fmt"
	"runtime"
	"strings"

	embed "github.com/Clinet/discordgo-embed" // go helper libary for parsing embeds
	"github.com/Zanos420/bcethabot/src/commands"
	customerror "github.com/Zanos420/bcethabot/src/error"
	"github.com/bwmarrin/discordgo"
)

type CmdHelp struct {
	prefix       string
	cmdInstances []commands.Command // all commands -> from cmd handler
}

func NewCmdHelp(prefix string, commands []commands.Command) *CmdHelp {
	return &CmdHelp{
		prefix:       prefix,
		cmdInstances: commands}
}

func (c *CmdHelp) Invokes() []string {
	return []string{"help", "h"} // Invokes and alias
}
func (c *CmdHelp) Description() string {
	return "Prints a help page"
}
func (c *CmdHelp) PermissionsNeeded() bool {
	return false
}

func (c *CmdHelp) CooldownLocked() bool {
	return true
}

func (c *CmdHelp) Exec(ctx *commands.Context) (err error) {
	if len(ctx.Args) > 0 {
		err = customerror.NewTooManyArgsError()
		return
	}

	emb := embed.NewEmbed()

	emb.SetTitle(fmt.Sprintf("%s - Helppage", ctx.Session.State.User.String()))
	emb.SetDescription(fmt.Sprintf("**__Commands - Bot Prefix:__ %s**", c.prefix))
	emb.SetThumbnail(ctx.Session.State.User.AvatarURL(""))
	emb.SetColor(0xE42D30)
	emb.InlineAllFields()

	for _, cmd := range c.cmdInstances {
		emb.AddField(cmd.Invokes()[0], fmt.Sprintf("%s (alias: %s)", cmd.Description(), strings.Join(cmd.Invokes()[:len(cmd.Invokes())], ", ")))
	}
	var owner *discordgo.User
	owner, err = ctx.Session.User("276414411881578497")
	if err != nil {
		return
	}
	emb.AddField("Go Version", runtime.Version())
	emb.AddField("Running on", runtime.GOOS)
	emb.SetFooter(fmt.Sprintf("Made by %s • If you find any bugs or want to contribute contact me", owner), owner.AvatarURL(""))

	_, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, emb.MessageEmbed)

	return
}
