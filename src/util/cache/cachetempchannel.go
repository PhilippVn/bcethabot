package cache

import (
	"fmt"
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
	Empty            bool
	Time_since_empty time.Time // time since the channel is empty
}

func NewCacheTempChannel() *CacheTempChannel {
	return &CacheTempChannel{
		Cache: &sync.Map{},
	}
}

func NewEntry(channel *discordgo.Channel, created_at time.Time, empty bool, time_since_empty time.Time) TmpEntry {
	return TmpEntry{
		Channel:          channel,
		Created_at:       created_at,
		Empty:            empty,
		Time_since_empty: time_since_empty,
	}
}

// Prints Cache. No Error Handling. Should only be used for debugging purposes
func (c *CacheTempChannel) PrintCache(s *discordgo.Session) {
	fmt.Println("--------------------\nTemp Channel Cache:")
	c.Cache.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(TmpEntry)

		usr, _ := s.User(k)
		created := v.Created_at.Format("15:04:05")
		empty_since := v.Time_since_empty.Format("15:04:05")
		fmt.Printf("-User:%s\n", usr.String())
		fmt.Printf("       -Channel    :%s\n", v.Channel.Name)
		fmt.Printf("       -Created    :%s\n", created)
		fmt.Printf("       -Empty      :%v\n", v.Empty)
		fmt.Printf("       -Empty Since: %s\n", empty_since)
		return true
	})
	fmt.Println("--------------------")
}
