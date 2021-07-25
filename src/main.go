package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"path/filepath"

	"github.com/Zanos420/bcethabot/src/commands"
	"github.com/Zanos420/bcethabot/src/commands/cmd"
	"github.com/Zanos420/bcethabot/src/events"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v3"
)

var (
	config  *Config // ptr to Config struct
	rootdir string  // path of the parent dir
)

/* config struct which caches all config data from the config file */
type Config struct {
	BOT struct {
		TOKEN  string `yaml:"TOKEN"`
		PREFIX string `yaml:"PREFIX"`
	} `yaml:"BOT"`
	VAR struct {
		CATEGORYID int `yaml:"CATEGORY_ID"`
		MODROLEID  int `yaml:"MOD_ROLE_ID"`
	} `yaml:"VAR"`
}

/* Config Parser that parses the config.yaml file to a Config struct */
func parseConfigFromYAMLFile(fileName string) (*Config, error) {
	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	var con *Config = new(Config)

	err = yaml.Unmarshal(buf, con)
	if err != nil {
		panic(err)
	}
	return con, err
}

/* initialized config data */
func init() {
	rootdir, err := os.Getwd()
	rootdir = filepath.Dir(rootdir)
	if err != nil {
		panic(err)
	}
	fmt.Println(rootdir)
	configPath, err := filepath.Abs(rootdir + "/config/config.yaml")
	if err != nil {
		panic(err)
	}
	fmt.Println(configPath)
	config, err = parseConfigFromYAMLFile(configPath)
	if err != nil {
		panic(err)
	}
}

/* bot creation */
func main() {
	// Create a new Discord session using the provided bot token.
	fmt.Printf("Running on Auth-Token: %s\n", config.BOT.TOKEN)
	bot, err := discordgo.New("Bot " + config.BOT.TOKEN)
	if err != nil {
		panic(err)
	}

	// bot should be able to listen to all events -> scaling shouldnt be a problem
	bot.Identify.Intents = discordgo.MakeIntent(
		discordgo.IntentsGuildMessages)
	// register events the bot is listening to
	registerEvents(bot)
	registerCommands(bot, config)

	err = bot.Open()
	if err != nil {
		panic(err)
	}

	fmt.Println("Bot is running... Press Ctlr-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	// Cleanly close down the Discord session.
	bot.Close()

}

func registerEvents(s *discordgo.Session) {
	s.AddHandler(events.NewMessageHandler().Handler)
	s.AddHandler(events.NewReadyHandler().Handler)
	fmt.Println("Successfully hooked all Event Handlers")
}

func registerCommands(s *discordgo.Session, cfg *Config) {
	cmdHandler := commands.NewCommandHandler(cfg.BOT.PREFIX)
	// Further Commands will be registered here
	cmdHandler.RegisterCommand(&cmd.CmdPing{})

	s.AddHandler(cmdHandler.HandleMessage)
	fmt.Println("Successfully hooked all Command Handlers")
}