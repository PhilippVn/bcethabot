# Black Cetha Manager Bot
Simple Bot with minimal configuration needed, written in Golang to support Black Cetha Network

# Configuration
## Creating Config File
To run the bot u need to create a `config.yaml` file inside the [config folder](bcethabot/config)
## Config File Format
The Config File has to look similar to this
```yaml
BOT:
 TOKEN: "<token>"
 PREFIX: "<prefix>"

VAR:
 CATEGORY_ID: <categoryid>
 MOD_ROLE_ID: <roleid>
```
## Explanation of Config File Values
- token       :   Token of your Discord Application/Bot. The Bot will use this token to auth with the websocket (Dont share your token!)
- prefix      :   Command prefix the bot is listening too (e.g. !help where '!' is the prefix and 'help' is the command issued)
- category_id :   Discord Channelcategory ID where the bot is going to create temporary channels on
- mod_role_id :   Roleid that is elidgeable to use the Bots Commands

# Running
To run the Bot after creating the config file you may simply use the provided `Makefile` which will create a binary depending on your os in [binaries](bcethabot/bin)
To compile the bot run `make` inside the root directory

# Specs and external dependencies
- Go `1.16`
- [discordgo](https://github.com/bwmarrin/discordgo) `0.23.2`
