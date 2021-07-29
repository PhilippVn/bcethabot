package cmd

//Nuke Command Module
import (
	"fmt"
	"sync"

	"github.com/Zanos420/bcethabot/src/commands"
	customerror "github.com/Zanos420/bcethabot/src/error"
	"github.com/Zanos420/bcethabot/src/error/internalerror"
	"github.com/Zanos420/bcethabot/src/util/cache"
	"github.com/bwmarrin/discordgo"
)

type CmdNuke struct {
	categoryID        string
	cacheTempChannels *cache.CacheTempChannel
	cacheOwners       *cache.CacheOwner
}

func NewCmdNuke(tempchannels *cache.CacheTempChannel, owners *cache.CacheOwner, categoryID string) *CmdNuke {
	return &CmdNuke{
		categoryID:        categoryID,
		cacheTempChannels: tempchannels,
		cacheOwners:       owners,
	}
}

func (c *CmdNuke) Invokes() []string {
	return []string{"nuke", "n"} // Invokes and alias
}
func (c *CmdNuke) Description() string {
	return "Deletes all Channels in the Temporary Channel Category + Cleans the Temp Channel Cache"
}
func (c *CmdNuke) PermissionsNeeded() bool {
	return true
}
func (c *CmdNuke) CooldownLocked() bool {
	return true
}
func (c *CmdNuke) Exec(ctx *commands.Context) (err error) {
	var tmpcategory *discordgo.Channel
	tmpcategory, err = ctx.Session.Channel(c.categoryID)
	if err != nil {
		err = customerror.NewInvalidConfigurationError()
		return
	}
	// Channel type is just an int so we can compare them like this
	if tmpcategory.Type != discordgo.ChannelTypeGuildCategory {
		err = customerror.NewInvalidConfigurationError()
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
		internalerror.Error(customerror.NewCustomError("Couldnt fetch Guild Channels"))
		err = nil // so it wont get send to User
		return
	}

	var del *discordgo.Channel
	for _, channel := range guildchannels {
		if channel.ParentID == c.categoryID {
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
	// clean caches
	c.cacheTempChannels.Cache = &sync.Map{}
	c.cacheOwners.Cache = &sync.Map{}
	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf(":white_check_mark: Nuked Temp Category: `%s`", tmpcategory.Name))
	return
}

func (c *CmdNuke) ExecDM(ctx *commands.Context) (err error) {
	var tmpcategory *discordgo.Channel
	tmpcategory, err = ctx.Session.Channel(c.categoryID)
	if err != nil {
		return
	}
	// Channel type is just an int so we can compare them like this
	if tmpcategory.Type != discordgo.ChannelTypeGuildCategory {
		err = customerror.NewInvalidConfigurationError()
		return
	}

	if len(ctx.Args) > 0 {
		err = customerror.NewTooManyArgsError()
		return
	}

	//nuke channels
	//collect all channels
	var guildchannels []*discordgo.Channel
	guildchannels, err = ctx.Session.GuildChannels(ctx.Session.State.Guilds[0].ID)

	if err != nil {
		return
	}

	var del *discordgo.Channel
	for _, channel := range guildchannels {
		if channel.ParentID == c.categoryID {
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
	// clean caches
	c.cacheTempChannels.Cache = &sync.Map{}
	c.cacheOwners.Cache = &sync.Map{}
	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf(":white_check_mark: Nuked Temp Category: `%s`", tmpcategory.Name))
	return
}
