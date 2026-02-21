package commands

import (
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/as7ar/artifex/embeds"
)

func PingCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!ping" {
		now := time.Now()
		messageTime := m.Timestamp
		ping := messageTime.Nanosecond() - now.Nanosecond()

		_, err := s.ChannelMessageSendEmbedReply(
			m.ChannelID,
			embeds.New().Color(0x7289da).Title("🏓 PONG!").Description("`"+strconv.Itoa(ping)+"`").Build(),
			m.Reference())
		if err != nil {
			return
		}
	}
}
