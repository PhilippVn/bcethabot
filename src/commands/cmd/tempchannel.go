package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Zanos420/bcethabot/src/commands"
	customerror "github.com/Zanos420/bcethabot/src/error"
	"github.com/bwmarrin/discordgo"
)

// This command lets you create temp channels for the given category
type CmdTempChannel struct {
	CATEGORY_ID       string
	TempChannels      sync.Map           // map userid -> Triple[channel, created_at,time_since_empty]
	TempChannelOwners sync.Map           // map channelid -> userid
	session           *discordgo.Session // session needed for deleting tmpchannels
}

type TempChannelTriple struct {
	Channel          discordgo.Channel
	Created_at       time.Time
	Time_since_empty time.Time // time since the channel is empty
}

func NewCmdTempChannel(categoryID string, bot *discordgo.Session) (ct *CmdTempChannel) {
	ct = &CmdTempChannel{
		CATEGORY_ID:       categoryID,
		TempChannels:      sync.Map{},
		TempChannelOwners: sync.Map{},
		session:           bot,
	}
	go ct.HeartBeatTempChannelDelete()
	return
}

// goroutine that will check the cached temp channels every minute and delete the empty ones that noone connected to for 1 min
func (c *CmdTempChannel) HeartBeatTempChannelDelete() {
	tmpcategory, err := c.session.Channel(c.CATEGORY_ID)
	if err != nil {
		panic(err)
	}
	if tmpcategory.Type != discordgo.ChannelTypeGuildCategory {
		fmt.Println("Category Channel not confugured correctly!")
		return
	}

	// will run every minute
	ticker := time.NewTicker(1 * time.Minute)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				// iterating over cached tmpChannels
				c.TempChannels.Range(func(key, value interface{}) bool {
					val := value.(TempChannelTriple)
					diff := time.Until(val.Time_since_empty)
					diffcreation := time.Until(val.Created_at)
					// channel can be empty for 2 min after creation and will be deleted after 1 min of being empty
					if diff.Minutes() <= -1 && diffcreation.Minutes() <= -2 {
						_, err := c.session.ChannelDelete(val.Channel.ID)
						if err != nil {
							fmt.Printf("Failed to heartbeat clean temporary Channel %s: %s\n", val.Channel.Name, err)
						}
						// delete keys from maps
						c.TempChannelOwners.Delete(val.Channel.ID)
						c.TempChannels.Delete(key)
					}
					return true
				})

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

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
func (c *CmdTempChannel) CooldownLocked() bool {
	return true
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

	t, ok := c.TempChannels.Load(ctx.Message.Author.ID)
	// user has created a temp channel already
	if ok {
		triple := t.(TempChannelTriple)
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf(":x: You already created a temp channel: `%s` %v seconds ago", triple.Channel.Name, (int)(-time.Until(triple.Created_at).Seconds())))
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
	// register channel creation in maps
	c.TempChannels.Store(ctx.Message.Author.ID, TempChannelTriple{*newchannel, time.Now(), time.Now()})
	c.TempChannelOwners.Store(newchannel.ID, ctx.Message.Author.ID)
	ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf(":white_check_mark: Temporary Channel `%s` created.\nAttention: Channel will be deleted when nobody is in channel!", chName))

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
