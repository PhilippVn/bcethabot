package mw

// Permission Middleware Module
import (
	"github.com/Zanos420/bcethabot/src/commands"
	customerror "github.com/Zanos420/bcethabot/src/error"
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
			err = customerror.NewMissingPermissionsError()
		}
	}()

	member, err := ctx.Session.GuildMember(ctx.Session.State.Guilds[0].ID, ctx.Message.Author.ID)

	if err != nil {
		err = customerror.NewNoSharedGuildError()
		return
	}

	// check if any of the authors roles is the mod role set in config.yaml
	for _, rID := range member.Roles {
		if mw.modRoleID == rID {
			next = true
			break
		}
	}
	return // if next is not set to true > excecution will stop
}

func (mw *MwPermissions) ExecDM(ctx *commands.Context, cmd commands.Command) (next bool, err error) {
	if !cmd.PermissionsNeeded() {
		next = true // deal with next mw/excecute command
		return
	}

	// will be excecuted after return of Exec()
	defer func() {
		if !next && err == nil {
			err = customerror.NewMissingPermissionsError()
		}
	}()

	member, err := ctx.Session.GuildMember(ctx.Session.State.Guilds[0].ID, ctx.Message.Author.ID)

	if err != nil {
		err = customerror.NewNoSharedGuildError()
		return
	}

	// check if any of the authors roles is the mod role set in config.yaml
	for _, rID := range member.Roles {
		if mw.modRoleID == rID {
			next = true
			break
		}
	}
	return // if next is not set to true > excecution will stop
}
