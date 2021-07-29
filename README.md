# Black Cetha Manager Bot
Simple Bot with minimal configuration needed, written in Golang to support Black Cetha Network

# Configuration
## Creating Config File
To run the bot u need to create a `config.yaml` file inside the [config folder](./config)
## Config File Format
The Config File has to look similar to this
```yaml
BOT:
 TOKEN: "<token>"
 PREFIX: "<prefix>"

VAR:
 CATEGORY_ID: "<categoryid>"
 MOD_CMD: "<true|false>"
 MOD_ROLE_ID: "<roleid>"
 CMD_COOLDOWN: "<true|false>"
 MOD_EXCLUDED: "<true|false>"
 CMD_COOLDOWN_SECONDS: "<seconds>"
```
## Explanation of Config File Values
- token               :   Token of your Discord Application/Bot. The Bot will use this token to auth with the websocket (Dont share your token!)
- prefix              :   Command prefix the bot is listening too (e.g. !help where '!' is the prefix and 'help' is the command issued)
- category_id         :   Discord Channelcategory ID where the bot is going to create temporary channels on
- mod_cmd             :   Boolean if the Bot should have mod only commands (default: true)
- mod_role_id         :   Roleid that is elidgeable to use the Bots Commands
- cmd_cooldown        :   Boolean if there should be a Command Cooldown (default: true)
- mod_excluded        :   Boolean if users with mod role should be excluded from cooldown  (default: true)
- cmd_cooldown_seconds:   Cooldown for each command per user in seconds

# Running
To run the Bot after creating the config file you may simply use the provided `Makefile` which will create a binary depending on your os in [binaries](./bin)
To compile the bot run `make` inside the root directory

# Specs and external dependencies
- Go `1.16`
- [discordgo Api Wrapper](https://github.com/bwmarrin/discordgo) `0.23.2`
- [discordgo Embed Libary](https://github.com/Clinet/discordgo-embed) `0.0.0`
- [yaml3](https://gopkg.in/yaml.v3) `3.0.0`
