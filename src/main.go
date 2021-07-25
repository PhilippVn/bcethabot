package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"path/filepath"

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
		fmt.Println("Fatal Error while creating Discord session: ", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	bot.AddHandler(messageCreate)

	// Just like the ping pong example, we only care about receiving message
	// events in this example.
	bot.Identify.Intents = discordgo.IntentsGuildMessages

	err = bot.Open()
	if err != nil {
		fmt.Println("Error while opening websocket connection: ", err)
		return
	}

	fmt.Println("Bot is running... Press Ctlr-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	// Cleanly close down the Discord session.
	bot.Close()

}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
//
// It is called whenever a message is created but only when it's sent through a
// server as we did not request IntentsDirectMessages.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong")
	}
	if m.Content == "huren" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "sohn")
	}
}
