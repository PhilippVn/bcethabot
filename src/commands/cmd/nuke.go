package cmd

import (
	"fmt"

	"github.com/Zanos420/bcethabot/src/commands"
	customerror "github.com/Zanos420/bcethabot/src/error"
	"github.com/bwmarrin/discordgo"
)

type CmdNuke struct {
	CATEGORY_ID string
}

func NewCmdNuke(categoryID string) *CmdNuke {
	return &CmdNuke{CATEGORY_ID: categoryID}
}

func (c *CmdNuke) Invokes() []string {
	return []string{"nuke", "n"} // Invokes and alias
}
func (c *CmdNuke) Description() string {
	return "Deletes all Channels in the Temporary Channel Category"
}
func (c *CmdNuke) PermissionsNeeded() bool {
	return true
}
func (c *CmdNuke) Exec(ctx *commands.Context) (err error) {
	var tmpcategory *discordgo.Channel
	tmpcategory, err = ctx.Session.Channel(c.CATEGORY_ID)
	if err != nil {
		return
	}
	// Channel type is just an int so we can compare them like this
	if tmpcategory.Type != discordgo.ChannelTypeGuildCategory {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Temp Channel category hasnt been configurated correctly!")
		return
	}

	if len(ctx.Args) > 0 {
		err = customerror.NewTooManyArgsError()
		return
	}

	//nuke channels
	//collect all channels
	var guildchannels []*discordgo.Channel
	guildchannels, err = ctx.Session.GuildChannels(ctx.Message.GuildID)

	if err != nil {
		return
	}

	var del *discordgo.Channel
	for _, channel := range guildchannels {
		if channel.ParentID == c.CATEGORY_ID {
			del, err = ctx.Session.ChannelDelete(channel.ID)
			if err != nil {
				return
			}
			_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf(":wastebasket: Deleted Temp Channel: `%s`", del.Name))
			if err != nil {
				return
			}
		}
	}

	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf(":white_check_mark: Nuked Temp Category: `%s`", tmpcategory.Name))
	return
}
