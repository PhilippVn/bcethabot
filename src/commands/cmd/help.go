package cmd

// Help Command Module
import (
	"fmt"
	"runtime"
	"strings"
	"time"

	embed "github.com/Clinet/discordgo-embed" // go helper libary for parsing embeds
	"github.com/Zanos420/bcethabot/src/commands"
	customerror "github.com/Zanos420/bcethabot/src/error"
	"github.com/bwmarrin/discordgo"
)

type CmdHelp struct {
	UpTime time.Time
}

func NewCmdHelp(uptime time.Time) *CmdHelp {
	return &CmdHelp{UpTime: uptime}
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
	emb.AddField("Uptime", formatDuration(time.Since(c.UpTime)))
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
	emb.AddField("Uptime", formatDuration(time.Since(c.UpTime)))
	emb.AddField("Go Version", runtime.Version())
	emb.AddField("Running on", runtime.GOOS)
	emb.SetFooter(fmt.Sprintf("Made by %s • If you find any bugs or want to contribute contact me", owner), owner.AvatarURL(""))

	_, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, emb.MessageEmbed)

	return
}

func formatDuration(duration time.Duration) (s string) {
	duration = duration.Round(time.Second) // pass by value
	DAY := 24 * time.Hour                  // libary doesnt provide it

	days := duration / DAY
	duration -= days * DAY
	hours := duration / time.Hour
	duration -= hours * time.Hour
	minutes := duration / time.Minute
	duration -= minutes * time.Minute
	seconds := duration / time.Second

	s = fmt.Sprintf("%0dd %0dh %0dm %0ds", days, hours, minutes, seconds)
	return
}
