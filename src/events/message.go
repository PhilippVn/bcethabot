package events

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type MessageHandler struct{}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{}
}

func (h *MessageHandler) Handler(s *discordgo.Session, e *discordgo.MessageCreate) {
	channel, err := s.Channel(e.ChannelID)
	if err != nil {
		fmt.Println("Failed to fetch Channel: ", err)
		return
	}

	fmt.Printf("User %s wrote in #%s: %s\n", e.Author, channel.Name, e.Content)
}
