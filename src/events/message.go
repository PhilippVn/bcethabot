package events

import (
	"fmt"
	"strings"

	customerror "github.com/Zanos420/bcethabot/src/error"
	"github.com/Zanos420/bcethabot/src/error/internalerror"
	"github.com/bwmarrin/discordgo"
)

type MessageHandler struct {
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{}
}

func (h *MessageHandler) Handler(s *discordgo.Session, e *discordgo.MessageCreate) {
	channel, err := s.Channel(e.ChannelID)
	if err != nil {
		internalerror.Error(customerror.NewCustomError("Failed to fetch Channel"))
		return
	}
	//fmt.Printf("User: %s, wrote %s, in: %s\n", e.Author.String(), e.Message.Content, channel.Name)
	if strings.Contains(e.Message.Content, "<@!739902368635813930>") {
		_, err := s.ChannelMessageSend(channel.ID, fmt.Sprintf("%s Ping mich noch einmal und es klatscht! :wave:", e.Author.Mention()))
		if err != nil {
			internalerror.Error(customerror.NewCustomError("Failed to sent message"))
		}
	}
}
