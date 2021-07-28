package mw

import (
	"github.com/Zanos420/bcethabot/src/commands"
	"github.com/bwmarrin/discordgo"
)

// permission middleware implements middleware
type MwPermissions struct {
	modRoleID string
}

func NewMwPermissions(modRoleID string) *MwPermissions {
	return &MwPermissions{
		modRoleID: modRoleID,
	}
}

func (mw *MwPermissions) Exec(ctx *commands.Context, cmd commands.Command) (next bool, err error) {
	// because this is a middleware it will be excecuted for every command -> filtering needed
	if !cmd.PermissionsNeeded() {
		next = true // deal with next mw/excecute command
		return
	}

	// will be excecuted after return of Exec()
	defer func() {
		if !next && err == nil {
			_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "You dont have permission to use this command!")
		}
	}()

	guild, err := ctx.Session.Guild(ctx.Message.GuildID)
	if err != nil {
		return
	}

	// if guild.OwnerID == ctx.Message.Author.ID { // owner of the guild
	// 	next = true
	// 	return
	// }

	roleMap := make(map[string]*discordgo.Role)

	// collect all guild roles
	for _, role := range guild.Roles {
		roleMap[role.ID] = role
	}

	// check if any of the authors roles has admin perms or is the mod role set in config.yaml
	for _, rID := range ctx.Message.Member.Roles {
		// if role, ok := roleMap[rID]; ok && role.Permissions&discordgo.PermissionAdministrator > 0 { // logical and with 0b1000 = admin
		// 	next = true
		// 	break
		// }
		if mw.modRoleID == rID {
			next = true
			break
		}
	}
	return // if next is not set to true > excecution will stop
}
