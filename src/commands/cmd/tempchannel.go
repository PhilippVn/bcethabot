package cmd

// Temp Channel Command Module
import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Zanos420/bcethabot/src/commands"
	customerror "github.com/Zanos420/bcethabot/src/error"
	"github.com/Zanos420/bcethabot/src/error/internalerror"
	"github.com/Zanos420/bcethabot/src/util/cache"
	"github.com/bwmarrin/discordgo"
)

// This command lets you create temp channels for the given category
type CmdTempChannel struct {
	categoryID        string
	cacheTempChannels *cache.CacheTempChannel
	cacheOwners       *cache.CacheOwner
	session           *discordgo.Session // session needed for deleting tmpchannels
}

func NewCmdTempChannel(tempChannels *cache.CacheTempChannel, owners *cache.CacheOwner, categoryID string, bot *discordgo.Session) (ct *CmdTempChannel) {
	ct = &CmdTempChannel{
		categoryID:        categoryID,
		cacheTempChannels: tempChannels,
		cacheOwners:       owners,
		session:           bot,
	}
	go ct.HeartBeatTempChannelDelete()
	return
}

// goroutine that will check the cached temp channels every minute and delete the empty ones that noone connected to for 1 min
func (c *CmdTempChannel) HeartBeatTempChannelDelete() {
	tmpcategory, err := c.session.Channel(c.categoryID)
	if err != nil {
		internalerror.Fatal(customerror.NewInvalidConfigurationError())
		return
	}
	if tmpcategory.Type != discordgo.ChannelTypeGuildCategory {
		internalerror.Fatal(customerror.NewInvalidConfigurationError())
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
				c.cacheTempChannels.Cache.Range(func(key, value interface{}) bool {
					val := value.(cache.TmpEntry)
					// if the channel isnt empty we dont want to delete it
					if !val.Empty {
						return true
					}
					diff := time.Until(val.Time_since_empty)
					diffcreation := time.Until(val.Created_at)
					// channel can be empty for 2 min after creation and will be deleted after 1 min of being empty
					if diff.Minutes() <= -1 && diffcreation.Minutes() <= -2 {
						_, err := c.session.ChannelDelete(val.Channel.ID)
						if err != nil {
							internalerror.Error(customerror.NewCustomError(fmt.Sprintf("Failed to heartbeat clean temporary Channel %s: %s\n", val.Channel.Name, err)))
						}
						// delete keys from maps
						c.cacheOwners.Cache.Delete(val.Channel.ID)
						c.cacheTempChannels.Cache.Delete(key)

						fmt.Println(">> AFTER HEARTBEATDELETE")
						c.cacheTempChannels.PrintCache(c.session)
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
	return []string{"tempchannel", "temp", "tmp", "t"} // Invokes and alias
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
	tmpcategory, err = ctx.Session.Channel(c.categoryID)
	if err != nil {
		return
	}
	// Channel type is just an int so we can compare them like this
	if tmpcategory.Type != discordgo.ChannelTypeGuildCategory {
		err = customerror.NewInvalidConfigurationError()
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

	if len([]rune(chName)) > 30 {
		err = customerror.NewCustomError("That name is a little bit to long. Please choose a shorter name (max 30 letters)!")
		return
	}

	if len([]rune(chName)) < 3 {
		err = customerror.NewCustomError("That name is a little bit to short. Please choose a longer name (min 3 letters)!")
		return
	}

	if limitSelected && (chLimit > 99 || chLimit < 1) {
		err = customerror.NewCustomError("Invalid User limit. Please choose a number between 1-99!")
		return
	}

	t, ok := c.cacheTempChannels.Cache.Load(ctx.Message.Author.ID)
	// user has created a temp channel already
	if ok {
		triple := t.(cache.TmpEntry)
		err = customerror.NewExistingTempChannelError(triple.Channel.Name, (int)(-time.Until(triple.Created_at).Seconds()))
		return
	}

	var tmpchannel discordgo.GuildChannelCreateData
	if limitSelected {
		tmpchannel = discordgo.GuildChannelCreateData{
			Name:      chName,
			Type:      discordgo.ChannelTypeGuildVoice,
			ParentID:  c.categoryID,
			UserLimit: chLimit,
		}
	} else {
		tmpchannel = discordgo.GuildChannelCreateData{
			Name:     chName,
			Type:     discordgo.ChannelTypeGuildVoice,
			ParentID: c.categoryID,
		}
	}
	newchannel, err := ctx.Session.GuildChannelCreateComplex(ctx.Message.GuildID, tmpchannel)
	if err != nil {
		return
	}
	// register channel creation in maps
	c.cacheTempChannels.Cache.Store(ctx.Message.Author.ID, cache.NewEntry(newchannel, time.Now(), true, time.Now()))
	c.cacheOwners.Cache.Store(newchannel.ID, ctx.Message.Author.ID)
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
	fmt.Println(">>AFTER TEMPCMD")
	c.cacheTempChannels.PrintCache(ctx.Session)
	return
}

func (c *CmdTempChannel) ExecDM(ctx *commands.Context) (err error) {
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

	if len([]rune(chName)) > 30 {
		err = customerror.NewCustomError("That name is a little bit to long. Please choose a shorter name (max 30 letters)!")
		return
	}

	if len([]rune(chName)) < 3 {
		err = customerror.NewCustomError("That name is a little bit to short. Please choose a longer name (min 3 letters)!")
		return
	}

	if limitSelected && (chLimit > 99 || chLimit < 1) {
		err = customerror.NewCustomError("Invalid User limit. Please choose a number between 1-99!")
		return
	}

	t, ok := c.cacheTempChannels.Cache.Load(ctx.Message.Author.ID)
	// user has created a temp channel already
	if ok {
		triple := t.(cache.TmpEntry)
		err = customerror.NewExistingTempChannelError(triple.Channel.Name, (int)(-time.Until(triple.Created_at).Seconds()))
		return
	}

	var tmpchannel discordgo.GuildChannelCreateData
	if limitSelected {
		tmpchannel = discordgo.GuildChannelCreateData{
			Name:      chName,
			Type:      discordgo.ChannelTypeGuildVoice,
			ParentID:  c.categoryID,
			UserLimit: chLimit,
		}
	} else {
		tmpchannel = discordgo.GuildChannelCreateData{
			Name:     chName,
			Type:     discordgo.ChannelTypeGuildVoice,
			ParentID: c.categoryID,
		}
	}
	newchannel, err := ctx.Session.GuildChannelCreateComplex(ctx.Session.State.Guilds[0].ID, tmpchannel)
	if err != nil {
		return
	}
	// register channel creation in maps
	c.cacheTempChannels.Cache.Store(ctx.Message.Author.ID, cache.NewEntry(newchannel, time.Now(), true, time.Now()))
	c.cacheOwners.Cache.Store(newchannel.ID, ctx.Message.Author.ID)
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
	fmt.Println(">>AFTER TEMPCMD")
	c.cacheTempChannels.PrintCache(ctx.Session)
	return
}
