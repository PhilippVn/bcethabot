package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Zanos420/bcethabot/src/commands"
	customerror "github.com/Zanos420/bcethabot/src/error"
	"github.com/bwmarrin/discordgo"
)

// This command lets you create temp channels for the given category
type CmdTempChannel struct {
	CATEGORY_ID string
}

func NewCmdTempChannel(categoryID string) *CmdTempChannel {
	return &CmdTempChannel{CATEGORY_ID: categoryID}
}

func (c *CmdTempChannel) Invokes() []string {
	return []string{"tempchannel", "temp", "tc"} // Invokes and alias
}
func (c *CmdTempChannel) Description() string {
	return "Creates a temporary Channel for you and your friends.\nUsage: <prefix>tempchannel [channelname] [optional:max useranzahl]"
}
func (c *CmdTempChannel) PermissionsNeeded() bool {
	return false
}
func (c *CmdTempChannel) Exec(ctx *commands.Context) (err error) {
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

	if len(ctx.Args) < 1 {
		err = customerror.NewNotEnoughArgsError()
		return
	}

	var chName string
	var chLimit int = 0
	var limitSelected bool = false
	if len(ctx.Args) == 1 {
		chName = ctx.Args[0]
	} else {
		//test if the user entered a number at the end -> user limit
		chLimit, err = strconv.Atoi(ctx.Args[len(ctx.Args)-1])
		if err == nil {
			limitSelected = true
		}
		err = nil
		if limitSelected {
			chName = strings.Join(ctx.Args[:len(ctx.Args)-1], " ")
		} else {
			chName = strings.Join(ctx.Args[:len(ctx.Args)], " ")
		}
	}

	if len([]rune(chName)) > 25 {
		err = customerror.NewCustomError("That name is a little bit to long. Please choose a shorter name (max 25 letters)!")
		return
	}

	if limitSelected && (chLimit > 99 || chLimit < 1) {
		err = customerror.NewCustomError("Invalid User limit. Please choose a number between 1-99!")
		return
	}

	var tmpchannel discordgo.GuildChannelCreateData
	if limitSelected {
		tmpchannel = discordgo.GuildChannelCreateData{
			Name:      chName,
			Type:      discordgo.ChannelTypeGuildVoice,
			ParentID:  c.CATEGORY_ID,
			UserLimit: chLimit,
		}
	} else {
		tmpchannel = discordgo.GuildChannelCreateData{
			Name:     chName,
			Type:     discordgo.ChannelTypeGuildVoice,
			ParentID: c.CATEGORY_ID,
		}
	}
	newchannel, err := ctx.Session.GuildChannelCreateComplex(ctx.Message.GuildID, tmpchannel)
	if err != nil {
		return
	}
	ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf(":white_check_mark: Temporärer Kanal `%s` erstellt.\nAchtung: Der Channel wird gelöscht sobald keiner mehr im Channel ist!", chName))

	var newchannelInv *discordgo.Invite = &discordgo.Invite{
		MaxAge:    120, //duration (in seconds) after which the invite expires
		MaxUses:   0,
		Temporary: false,
	}

	newchannelInv, err = ctx.Session.ChannelInviteCreate(newchannel.ID, *newchannelInv)
	if err != nil {
		return
	}
	ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf("https://discord.gg/%s", newchannelInv.Code))

	return
}
