package commands

// Context Module
import "github.com/bwmarrin/discordgo"

// Context Object
// Every Command issued has a context with the needed information
// [prefix][invoke/alias] [1st arg] [2nd arg] [3rd arg]...
type Context struct {
	Session *discordgo.Session // discord session that rec
	Message *discordgo.Message // Message Object that issues the command
	Args    []string           // A Command can have arguments, e.g. <prefix>ban <membername> <duration> <reason>
	Handler *CommandHandler    // Command Handler (e.g. usefull for db access)
}
