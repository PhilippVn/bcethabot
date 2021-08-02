package events

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// struct for database integrationaccess later
type ReadyHandler struct {
	Prefix string
}

// constructor for a Ready Handler
func NewReadyHandler(prefix string) *ReadyHandler {
	return &ReadyHandler{Prefix: prefix}
}

// Handler Method of "Class" Ready Handler
// notice the type of the Handler method which is a method of ReadyHandler
func (h *ReadyHandler) Handler(s *discordgo.Session, e *discordgo.Ready) {
	fmt.Println("----------------------------------------")
	fmt.Println("Bot session is ready")
	fmt.Printf("Logged in as %s\n", e.User)
	fmt.Println("----------------------------------------")
	// Update Presence
	s.UpdateListeningStatus(fmt.Sprintf("Prefix %s", h.Prefix))
}
