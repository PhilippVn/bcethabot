package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"path/filepath"

	"github.com/Zanos420/bcethabot/src/commands"
	"github.com/Zanos420/bcethabot/src/commands/cmd"
	"github.com/Zanos420/bcethabot/src/commands/mw"
	"github.com/Zanos420/bcethabot/src/error/internalerror"
	"github.com/Zanos420/bcethabot/src/events"
	"github.com/Zanos420/bcethabot/src/util/cache"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v3"
)

var (
	config      *Config
	rootdirPath string

	cacheTempChannels *cache.CacheTempChannel
	cacheOwners       *cache.CacheOwner
)

/* config struct which caches all config data from the config file */
type Config struct {
	BOT struct {
		TOKEN  string `yaml:"TOKEN"`
		PREFIX string `yaml:"PREFIX"`
	} `yaml:"BOT"`
	VAR struct {
		CATEGORYID      string `yaml:"CATEGORY_ID"`
		MODCMD          string `yaml:"MOD_CMD"`
		MODROLEID       string `yaml:"MOD_ROLE_ID"`
		COOLDOWN        string `yaml:"CMD_COOLDOWN"`
		MODEXCLUDED     string `yaml:"MOD_EXCLUDED"`
		COOLDOWNSECONDS string `yaml:"CMD_COOLDOWN_SECONDS"`
	} `yaml:"VAR"`
}

/* Config Parser that parses the config.yaml file to a Config struct */
func parseConfigFromYAMLFile(fileName string) (*Config, error) {
	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		internalerror.Fatal(err)
		os.Exit(0)
	}
	var con *Config = new(Config)

	err = yaml.Unmarshal(buf, con)
	if err != nil {
		internalerror.Fatal(err)
		os.Exit(0)
	}
	return con, err
}

/* initialized config data */
func init() {
	rootdir, err := os.Getwd()
	rootdir = filepath.Dir(rootdir)
	if err != nil {
		internalerror.Fatal(err)
		os.Exit(0)
	}
	fmt.Println(rootdir)
	configPath, err := filepath.Abs(rootdir + "/config/config.yaml")
	if err != nil {
		internalerror.Fatal(err)
		os.Exit(0)
	}
	fmt.Println(configPath)
	config, err = parseConfigFromYAMLFile(configPath)
	if err != nil {
		internalerror.Fatal(err)
		os.Exit(0)
	}
}

/* bot creation */
func main() {
	// Create a new Discord session using the provided bot token.
	internalerror.Info("Running on Auth-Token: %s", config.BOT.TOKEN)
	bot, err := discordgo.New("Bot " + config.BOT.TOKEN)
	if err != nil {
		internalerror.Fatal(err)
		os.Exit(0)
	}

	// bot should be able to listen to all events -> scaling shouldnt be a problem
	bot.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged | discordgo.IntentsGuildMembers)
	bot.StateEnabled = true // cache guild, members, ...
	//initialize internal caches (temporary channel functions)
	initializeCaches()
	// register events the bot is listening to
	registerCommands(bot, config)
	// Important: register Events after the commands! -> tmpcmd has to be initialized for voicestateupdate event
	registerEvents(bot, config)

	err = bot.Open()
	if err != nil {
		internalerror.Fatal(err)
		os.Exit(0)
	}

	fmt.Println("Bot is running... Press Ctlr-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	// Cleanly close down the Discord session.
	bot.Close()

}

func initializeCaches() {
	cacheTempChannels = cache.NewCacheTempChannel()
	cacheOwners = cache.NewCacheOwner()
}

func registerEvents(s *discordgo.Session, cfg *Config) {
	s.AddHandler(events.NewMessageHandler().Handler)
	s.AddHandler(events.NewReadyHandler().Handler)
	s.AddHandler(events.NewVoiceStateUpdateHandler(cacheTempChannels, cacheOwners, cfg.VAR.CATEGORYID).Handler)
	internalerror.Info("Successfully hooked all Event Handlers")
}

func registerCommands(s *discordgo.Session, cfg *Config) {
	// init command handler
	cmdHandler := commands.NewCommandHandler(cfg.BOT.PREFIX)

	// Commands
	cmdHandler.RegisterCommand(cmd.NewCmdPing())
	cmdTempHandler := cmd.NewCmdTempChannel(cacheTempChannels, cacheOwners, cfg.VAR.CATEGORYID, s) // save reference for help command which depends on the cmdList
	cmdHandler.RegisterCommand(cmdTempHandler)
	cmdHandler.RegisterCommand(cmd.NewCmdNuke(cacheTempChannels, cacheOwners, cfg.VAR.CATEGORYID))

	// Help command after registrating all other commands
	cmdHandler.RegisterCommand(cmd.NewCmdHelp())

	// Middlewares
	modcmd, err := strconv.ParseBool(cfg.VAR.MODCMD)
	if err != nil {
		internalerror.Fatal(err)
		os.Exit(0)
	}
	if modcmd {
		cmdHandler.RegisterMiddleware(mw.NewMwPermissions(cfg.VAR.MODROLEID))
	}

	cooldown, err := strconv.ParseBool(cfg.VAR.COOLDOWN)
	if err != nil {
		internalerror.Fatal(err)
		os.Exit(0)
	}
	modexcluded, err := strconv.ParseBool(cfg.VAR.MODEXCLUDED)
	if err != nil {
		internalerror.Fatal(err)
		os.Exit(0)
	}
	if cooldown {
		cooldownseconds, err := strconv.Atoi(cfg.VAR.COOLDOWNSECONDS)
		if err != nil {
			internalerror.Fatal(err)
			os.Exit(0)
		}
		cmdHandler.RegisterMiddleware(mw.NewMwCooldown(cooldownseconds, modexcluded, cfg.VAR.MODROLEID))
	}

	s.AddHandler(cmdHandler.HandleMessage)
	internalerror.Info("Successfully hooked all Command Handlers and Middlewares")

}
