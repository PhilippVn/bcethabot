package events

import (
	"fmt"
	"time"

	"github.com/Zanos420/bcethabot/src/util/cache"
	"github.com/bwmarrin/discordgo"
)

// this event is used to update the empty_since time for the cmdtmpchannel
type VoiceStateUpdateHandler struct {
	cacheTempChannels *cache.CacheTempChannel
	cacheOwners       *cache.CacheOwner
	categoryID        string
}

// Needs the Map from cmdtmpchannel
func NewVoiceStateUpdateHandler(tempChannels *cache.CacheTempChannel, owners *cache.CacheOwner, categoryId string) *VoiceStateUpdateHandler {
	return &VoiceStateUpdateHandler{
		cacheTempChannels: tempChannels,
		cacheOwners:       owners,
		categoryID:        categoryId,
	}
}

func (h *VoiceStateUpdateHandler) Handler(s *discordgo.Session, e *discordgo.VoiceStateUpdate) {
	// VoiceStateUpdate has 1 or two states
	// VoiceState and Before Update
	// Before update only available when previous VoiceState was cached

	before := e.BeforeUpdate
	now := e.VoiceState
	guild, err := s.State.Guild(e.GuildID)
	if err != nil {
		fmt.Println("Failed to fetch Guild from State: ", err)
		return
	}

	if before == nil && now == nil {
		fmt.Println("Fatal Error while fetching VoiceStates: Both nil values")
		return
	}

	// user either switched or disconnected muted, got muted, etc

	// check if channel is a temp channel
	if before == nil || before.ChannelID == "" {
		// we only have to check now -> if before doesnt exist it means this is a first time connect from this user to the now channel
		ch, err := s.Channel(now.ChannelID)
		if err != nil {
			fmt.Println("Failed to fetch temp channel: ", err)
			return
		}
		// not a temp channel
		if ch.ParentID != h.categoryID {
			fmt.Println("Not a Temp Channel")
			return
		}
		//temp channel -> count users connected
		connectedUsers := countConnectedUsers(guild.VoiceStates, ch.ID)
		// Updating the Cache to set the tmp channel Empty = true or false
		// to get the old values from cache we have to find the key (ownerId of the channel/person who created channel via command)
		chO, ok := h.cacheOwners.Cache.Load(ch.ID)
		if !ok {
			// channel not cached -> not created by bot
			fmt.Println("Error while trying to fetch ownerid of tempchannel to update empty_since time of tmp channel")
			return
		}
		chOwnerID := chO.(string)
		mapE, ok := h.cacheTempChannels.Cache.Load(chOwnerID)
		if !ok {
			// channel not cached -> not created by bot
			fmt.Println("Error while trying to fetch ownerid of tempchannel to update empty_since time of tmp channel")
			return
		}
		mapEntryOld := mapE.(cache.TmpEntry)
		var mapEntryNew cache.TmpEntry
		if connectedUsers == 0 {
			// Set Empty=true
			mapEntryNew = cache.NewEntry(mapEntryOld.Channel, mapEntryOld.Created_at, true, time.Now())
		} else {
			// Set Empty=false
			mapEntryNew = cache.NewEntry(mapEntryOld.Channel, mapEntryOld.Created_at, false, mapEntryOld.Time_since_empty)
		}
		h.cacheTempChannels.Cache.Store(chOwnerID, mapEntryNew)
	} else if now == nil || now.ChannelID == "" {
		// we only have to check before -> if now doesnt exist this means the user disconnected completely
		ch, err := s.Channel(before.ChannelID)
		if err != nil {
			fmt.Println("Failed to fetch temp channel: ", err)
			return
		}
		// not a temp channel
		if ch.ParentID != h.categoryID {
			fmt.Println("Not a Temp Channel")
			return
		}
		//temp channel -> count users connected
		connectedUsers := countConnectedUsers(guild.VoiceStates, ch.ID)
		chO, ok := h.cacheOwners.Cache.Load(ch.ID)
		if !ok {
			fmt.Println("Error while trying to fetch ownerid of tempchannel to update empty_since time of tmp channel")
			return
		}
		chOwnerID := chO.(string)
		mapE, ok := h.cacheTempChannels.Cache.Load(chOwnerID)
		if !ok {
			fmt.Println("Error while trying to fetch ownerid of tempchannel to update empty_since time of tmp channel")
			return
		}
		mapEntryOld := mapE.(cache.TmpEntry)
		var mapEntryNew cache.TmpEntry
		if connectedUsers == 0 {
			mapEntryNew = cache.NewEntry(mapEntryOld.Channel, mapEntryOld.Created_at, true, time.Now())
		} else {
			mapEntryNew = cache.NewEntry(mapEntryOld.Channel, mapEntryOld.Created_at, false, mapEntryOld.Time_since_empty)
		}
		h.cacheTempChannels.Cache.Store(chOwnerID, mapEntryNew)

	} else {
		var channelBeforeIsTemp bool = false
		var channelNowIsTemp bool = false
		// there is a before and after state
		channelBefore, err := s.Channel(before.ChannelID)
		if err != nil {
			fmt.Println("Failed to fetch temp channel from BeforeState: ", err)
			return
		}
		channelNow, err := s.Channel(now.ChannelID)
		fmt.Println(now.ChannelID)
		fmt.Println(e.VoiceState.ChannelID)
		if err != nil {
			fmt.Println("Failed to fetch temp channel from NowState: ", err)
			return
		}

		if channelBefore.ParentID == h.categoryID {
			channelBeforeIsTemp = true
		}
		if channelNow.ParentID == h.categoryID {
			channelNowIsTemp = true
		}

		if channelBefore.ID == channelNow.ID {
			//both the same channel (no switch)
			if !channelNowIsTemp {
				fmt.Println("Not a Temp Channel")
				return
			}
			// temp channel
			connectedUsers := countConnectedUsers(guild.VoiceStates, channelNow.ID)
			chO, ok := h.cacheOwners.Cache.Load(channelNow.ID)
			if !ok {
				// channel not cached -> not created by bot
				fmt.Println("Error while trying to fetch ownerid of tempchannel to update empty_since time of tmp channel")
				return
			}
			chOwnerID := chO.(string)
			mapE, ok := h.cacheTempChannels.Cache.Load(chOwnerID)
			if !ok {
				// channel not cached -> not created by bot
				fmt.Println("Error while trying to fetch ownerid of tempchannel to update empty_since time of tmp channel")
				return
			}
			mapEntryOld := mapE.(cache.TmpEntry)
			var mapEntryNew cache.TmpEntry
			if connectedUsers == 0 {
				// Set Empty=true
				mapEntryNew = cache.NewEntry(mapEntryOld.Channel, mapEntryOld.Created_at, true, time.Now())
			} else {
				// Set Empty=false
				mapEntryNew = cache.NewEntry(mapEntryOld.Channel, mapEntryOld.Created_at, false, mapEntryOld.Time_since_empty)
			}
			h.cacheTempChannels.Cache.Store(chOwnerID, mapEntryNew)
		} else {
			// update both
			if channelBeforeIsTemp {
				connectedUsers := countConnectedUsers(guild.VoiceStates, channelBefore.ID)
				chO, ok := h.cacheOwners.Cache.Load(channelBefore.ID)
				if !ok {
					// channel not cached -> not created by bot
					fmt.Println("Error while trying to fetch ownerid of tempchannel to update empty_since time of tmp channel")
					return
				}
				chOwnerID := chO.(string)
				mapE, ok := h.cacheTempChannels.Cache.Load(chOwnerID)
				if !ok {
					// channel not cached -> not created by bot
					fmt.Println("Error while trying to fetch ownerid of tempchannel to update empty_since time of tmp channel")
					return
				}
				mapEntryOld := mapE.(cache.TmpEntry)
				var mapEntryNew cache.TmpEntry
				if connectedUsers == 0 {
					// Set Empty=true
					mapEntryNew = cache.NewEntry(mapEntryOld.Channel, mapEntryOld.Created_at, true, time.Now())
				} else {
					// Set Empty=false
					mapEntryNew = cache.NewEntry(mapEntryOld.Channel, mapEntryOld.Created_at, false, mapEntryOld.Time_since_empty)
				}
				h.cacheTempChannels.Cache.Store(chOwnerID, mapEntryNew)
			}
			if channelNowIsTemp {
				connectedUsers := countConnectedUsers(guild.VoiceStates, channelNow.ID)
				chO, ok := h.cacheOwners.Cache.Load(channelNow.ID)
				if !ok {
					// channel not cached -> not created by bot
					fmt.Println("Error while trying to fetch ownerid of tempchannel to update empty_since time of tmp channel")
					return
				}
				chOwnerID := chO.(string)
				mapE, ok := h.cacheTempChannels.Cache.Load(chOwnerID)
				if !ok {
					// channel not cached -> not created by bot
					fmt.Println("Error while trying to fetch ownerid of tempchannel to update empty_since time of tmp channel")
					return
				}
				mapEntryOld := mapE.(cache.TmpEntry)
				var mapEntryNew cache.TmpEntry
				if connectedUsers == 0 {
					// Set Empty=true
					mapEntryNew = cache.NewEntry(mapEntryOld.Channel, mapEntryOld.Created_at, true, time.Now())
				} else {
					// Set Empty=false
					mapEntryNew = cache.NewEntry(mapEntryOld.Channel, mapEntryOld.Created_at, false, mapEntryOld.Time_since_empty)
				}
				h.cacheTempChannels.Cache.Store(chOwnerID, mapEntryNew)
			}
		}

	}
	fmt.Println(">> AFTER VOICESTATEEVT")
	h.cacheTempChannels.PrintCache(s)

}

// count the connected users to a voicechannel
func countConnectedUsers(guildVoiceStates []*discordgo.VoiceState, channelID string) int {
	var connectedUsers int = 0
	for _, voiceState := range guildVoiceStates {
		if voiceState.ChannelID == channelID {
			connectedUsers++
		}
	}
	return connectedUsers
}

// prints a VoiceState. No Error Handlung should only be used for debugging
func printVoiceState(s *discordgo.Session, vs *discordgo.VoiceState) {
	if vs == nil {
		fmt.Println("Voice State: nil")
		return
	}
	ch, err := s.Channel(vs.ChannelID)
	if err != nil {
		return
	}
	usr, err := s.User(vs.UserID)
	if err != nil {
		return
	}
	g, err := s.Guild(vs.GuildID)
	if err != nil {
		return
	}
	fmt.Println("Voice State:")
	fmt.Printf("- Guild:%s\n", g.Name)
	fmt.Printf("- Channel:%s\n", ch.Name)
	fmt.Printf("- User:%s\n", usr.String())
	fmt.Printf("- Mute:%v\n", vs.Mute)
	fmt.Printf("- SelfMute:%v\n", vs.SelfMute)
	fmt.Printf("- Deaf:%v\n", vs.Deaf)
	fmt.Printf("- SelfMute:%v\n", vs.SelfDeaf)
	fmt.Printf("- Supress:%v\n", vs.Suppress)
}
