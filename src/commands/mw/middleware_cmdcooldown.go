package mw

import (
	"fmt"
	"sync"
	"time"

	"github.com/Zanos420/bcethabot/src/commands"
	customerror "github.com/Zanos420/bcethabot/src/error"
	"github.com/bwmarrin/discordgo"
)

// All Users without special role get cooldown blocked when they issue a cooldown command
// works with FiFo principle
type MwCooldown struct {
	usersOnCooldown sync.Map // map of userirds to lists of Pair{cmd,issued_at}
	coolDownSeconds int      // seconds of cooldown for each command
	modexcluded     bool     // wether mods are excluded from the cooldown
	modroleID       string   // mod role
}

type Pair struct {
	cmd       commands.Command // command that the user issued
	issued_at time.Time        // time the command has been issued
}

func NewMwCooldown(cooldownInSeconds int, modExcluded bool, modroleId string) (mw *MwCooldown) {
	mw = &MwCooldown{
		usersOnCooldown: sync.Map{},
		coolDownSeconds: cooldownInSeconds,
		modexcluded:     modExcluded,
		modroleID:       modroleId,
	}
	go mw.HeartBeatCleanUp()
	return
}

//task that runs every second and removes active command cooldowns
func (mw *MwCooldown) HeartBeatCleanUp() {
	// will run every second
	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				//fmt.Println("Cleaning cooldowns...")
				var removedCooldowns int
				// iterate over all userids and ther slices of Pairs
				mw.usersOnCooldown.Range(func(key, value interface{}) bool {
					pairList := value.([]Pair)
					// if there are no cmd on cooldown for this uid we just delete the entry
					if len(pairList) == 0 {
						mw.usersOnCooldown.Delete(key.(string))
					} else {
						// check if the oldest issued command (first in the Slice of Pairs) is older than 10 seconds
						diff := time.Until(pairList[0].issued_at)
						//fmt.Printf("Uid: %s, cmd: %s, time diff in seconds: %v\n", key.(string), pairList[0].cmd.Invokes()[0], diff.Seconds())
						if diff.Seconds() <= -float64(mw.coolDownSeconds) {
							// if so we delete the pair of cmd and issued_at
							value = pairList[1:]
							mw.usersOnCooldown.Store(key.(string), value)
							removedCooldowns++
						}
					}
					return true
				})
				//fmt.Printf("Removed %v cooldowns\n", removedCooldowns)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (mw *MwCooldown) Exec(ctx *commands.Context, cmd commands.Command) (next bool, err error) {
	// because this is a middleware it will be excecuted for every command -> filtering needed
	if !cmd.CooldownLocked() {
		next = true // deal with next mw/excecute command
		return
	}

	// if mods are excluded from the cooldown
	if mw.modexcluded {
		var guild *discordgo.Guild
		guild, err = ctx.Session.Guild(ctx.Message.GuildID)
		if err != nil {
			return
		}

		roleMap := make(map[string]*discordgo.Role)

		// collect all guild roles
		for _, role := range guild.Roles {
			roleMap[role.ID] = role
		}

		// check if any of the authors roles has admin perms or is the mod role set in config.yaml
		for _, rID := range ctx.Message.Member.Roles {
			if mw.modroleID == rID {
				next = true
				return
			}
		}
	}

	cmdL, ok := mw.usersOnCooldown.Load(ctx.Message.Author.ID)
	if ok {
		var pairList []Pair = cmdL.([]Pair)
		// user already on cooldown
		// check cmd list the user is on cooldown for
		for _, pair := range pairList {
			if pair.cmd == cmd {
				err = customerror.NewCommandOnCooldownErrorError((int)(float64(mw.coolDownSeconds) + time.Until(pairList[0].issued_at).Seconds()))
				//mw.printMap(mw.usersOnCooldown)
				return
			}
		}
		// add this command to his cooldown list -> first item in the list is always the longest command being issued
		pairList = append(pairList, Pair{cmd: cmd, issued_at: time.Now()})
		mw.usersOnCooldown.Store(ctx.Message.Author.ID, pairList)
		next = true
		//mw.printMap(mw.usersOnCooldown)
		return
	} else {
		// create new entry with the current command
		var pairList []Pair = []Pair{}
		pairList = append(pairList, Pair{cmd: cmd, issued_at: time.Now()})
		mw.usersOnCooldown.Store(ctx.Message.Author.ID, pairList)
		next = true
		//mw.printMap(mw.usersOnCooldown)
		return
	}
}

func PrintMap(mapping sync.Map) {
	fmt.Println("Map:{")
	mapping.Range(func(key, value interface{}) bool {
		fmt.Printf("Userid:%s\n", key.(string))
		for _, pair := range value.([]Pair) {
			fmt.Printf("--- command: %s - %v\n", pair.cmd.Invokes()[0], pair.issued_at)
		}
		return true
	})
	fmt.Printf("}\n")
}
