package events

import (
	"fmt"
	"sync"
	"time"

	"github.com/Zanos420/bcethabot/src/commands/cmd"
	"github.com/bwmarrin/discordgo"
)

// this event is used to update the empty_since time for the cmdtmpchannel

// struct for database integrationaccess later
type VoiceStateUpdateHandler struct {
	tmpChannelsMapImport       *sync.Map
	tmpChannelsOwnersMapImport *sync.Map
	categoryID                 string
}

// Needs the Map from cmdtmpchannel
func NewVoiceStateUpdateHandler(tmpChannelsMap *sync.Map, tmpChannelsOwnersMap *sync.Map, categoryId string) *VoiceStateUpdateHandler {
	return &VoiceStateUpdateHandler{
		tmpChannelsMapImport:       tmpChannelsMap,
		tmpChannelsOwnersMapImport: tmpChannelsOwnersMap,
		categoryID:                 categoryId,
	}
}

func (h *VoiceStateUpdateHandler) Handler(s *discordgo.Session, e *discordgo.VoiceStateUpdate) {

	channelIDNew := e.VoiceState.ChannelID
	var channelIDOld string

	var chID string // will be the final channel id

	if e.BeforeUpdate == nil {

		// user connected to voice channel -> channel cant be empty
		chID = channelIDNew
	} else {
		channelIDOld = e.BeforeUpdate.ChannelID
		chID = channelIDOld
	}

	// we want to ignore all these Voice States Changes
	if e.VoiceState.Mute || e.VoiceState.SelfMute || e.VoiceState.Deaf || e.VoiceState.SelfDeaf || e.VoiceState.Suppress {
		return
	}
	if e.BeforeUpdate != nil {
		if e.BeforeUpdate.Mute || e.BeforeUpdate.SelfMute || e.BeforeUpdate.Deaf || e.BeforeUpdate.SelfDeaf || e.BeforeUpdate.Suppress {
			return
		}
	}

	// user either switched or disconnected

	// check if channel is a temp channel
	channel, err := s.Channel(chID)
	if err != nil {

		fmt.Println("Failed to fetch temp channel: ", err)
		return
	}

	// not a temp channel
	if channel.ParentID != h.categoryID {

		return
	}

	// temp channel -> is it empty now?
	guild, err := s.Guild(e.GuildID)

	if err != nil {

		fmt.Println("Failed to fetch Guild: ", err)
		return
	}

	allVoiceStates := guild.VoiceStates
	var connectedUsers int = 0

	for _, voiceState := range allVoiceStates {
		if voiceState.ChannelID == chID {
			connectedUsers++
		}
	}

	//channel is empty
	if connectedUsers == 0 {

		//c.TempChannels.Store(ctx.Message.Author.ID, TempChannelTriple{*newchannel, time.Now(), time.Now()})
		// we have to find the key (person who issued the command to create the channel)
		chO, ok := h.tmpChannelsOwnersMapImport.Load(chID) // will give us the owner
		if !ok {
			// channel not cached -> not created by bot
			fmt.Println("Error while trying to fetch ownerid of tempchannel to update empty_since time of tmp channel")
			return
		}
		chOwnerID := chO.(string)

		// map userid -> Triple[channel, created_at,time_since_empty]
		mapT, ok := h.tmpChannelsMapImport.Load(chOwnerID)
		if !ok {
			// channel not cached -> not created by bot
			fmt.Println("Error while trying to fetch ownerid of tempchannel to update empty_since time of tmp channel")
			return
		}

		mapTripleOld := mapT.(cmd.TempChannelTriple)

		mapTripleNew := cmd.TempChannelTriple{
			Channel:          mapTripleOld.Channel,
			Created_at:       mapTripleOld.Created_at,
			Time_since_empty: time.Now(),
		}

		h.tmpChannelsMapImport.Store(chOwnerID, mapTripleNew)

	}
}
