package cache

import (
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

// This is a cache for temporary channels
type CacheTempChannel struct {
	Cache *sync.Map // Map [userid] -> CacheTempChannelTriple{channelptr, created_at, time_since_empty}
}

// values of the Cache
type TmpEntry struct {
	Channel          *discordgo.Channel
	Created_at       time.Time
	Time_since_empty time.Time // time since the channel is empty
}

func (c *CacheTempChannel) NewCacheTempChannel() *CacheTempChannel {
	return &CacheTempChannel{
		Cache: &sync.Map{},
	}
}

func NewEntry(channel *discordgo.Channel, created_at time.Time, time_since_empty time.Time) TmpEntry {
	return TmpEntry{
		Channel:          channel,
		Created_at:       created_at,
		Time_since_empty: time_since_empty,
	}
}
