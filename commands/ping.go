package commands

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func PingCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!ping" {
		now := time.Now()
		messageTime := m.Timestamp
		ping := messageTime.Nanosecond() - now.Nanosecond()

		_, err := s.ChannelMessageSendEmbedReply(m.ChannelID, embed.New, m.Reference())
		if err != nil {
			return
		}
	}
}
