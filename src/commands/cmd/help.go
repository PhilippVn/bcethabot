package cmd

// Help Command Module
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
}

func NewCmdHelp() *CmdHelp {
	return &CmdHelp{}
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
	emb.SetDescription(fmt.Sprintf("**__Commands - Bot Prefix:__ %s**\nCommands with a * require special permissions", ctx.Handler.Prefix))
	emb.SetThumbnail(ctx.Session.State.User.AvatarURL(""))
	emb.SetColor(0xE42D30)
	emb.InlineAllFields()

	for _, cmd := range ctx.Handler.CmdInstances {
		if cmd.PermissionsNeeded() {
			emb.AddField(fmt.Sprint("*", cmd.Invokes()[0]), fmt.Sprintf("%s (alias: %s)", cmd.Description(), strings.Join(cmd.Invokes()[:len(cmd.Invokes())], ", ")))
		} else {
			emb.AddField(cmd.Invokes()[0], fmt.Sprintf("%s (alias: %s)", cmd.Description(), strings.Join(cmd.Invokes()[:len(cmd.Invokes())], ", ")))
		}

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

func (c *CmdHelp) ExecDM(ctx *commands.Context) (err error) {
	if len(ctx.Args) > 0 {
		err = customerror.NewTooManyArgsError()
		return
	}
	emb := embed.NewEmbed()

	emb.SetTitle(fmt.Sprintf("%s - Helppage", ctx.Session.State.User.String()))
	emb.SetDescription(fmt.Sprintf("**__Commands - Bot Prefix:__ %s**\nCommands with a * require special permissions", ctx.Handler.Prefix))
	emb.SetThumbnail(ctx.Session.State.User.AvatarURL(""))
	emb.SetColor(0xE42D30)
	emb.InlineAllFields()

	for _, cmd := range ctx.Handler.CmdInstances {
		if cmd.PermissionsNeeded() {
			emb.AddField(fmt.Sprint("*", cmd.Invokes()[0]), fmt.Sprintf("%s (alias: %s)", cmd.Description(), strings.Join(cmd.Invokes()[:len(cmd.Invokes())], ", ")))
		} else {
			emb.AddField(cmd.Invokes()[0], fmt.Sprintf("%s (alias: %s)", cmd.Description(), strings.Join(cmd.Invokes()[:len(cmd.Invokes())], ", ")))
		}

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
