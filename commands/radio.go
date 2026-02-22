package commands

import (
	"context"
	"fmt"

	"github.com/as7ar/noori/config"
	"github.com/as7ar/noori/embeds"
	"github.com/as7ar/noori/logger"
	"github.com/as7ar/noori/utils"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/state/store/defaultstore"
)

func RadioCommand(s *session.Session, c *gateway.MessageCreateEvent) {
	guildID, chatID := c.GuildID, c.ChannelID
	vscID, _ := utils.GetUserVoiceChannelID(state.NewFromSession(s, defaultstore.New()), guildID, c.Author.ID)

	logger.Debug("GuildID:", guildID, "ChatID:", chatID)
	logger.Debug("VSC:", vscID)

	if vscID == 0 {
		_, err := s.SendEmbedReply(
			c.ChannelID, c.Message.ID,
			embeds.New().
				Color(config.SymbolColor).
				Title("🚨 Warning").
				Description("You must be entered in `voice chat`").Build())
		if err != nil {
			fmt.Println("send error:", err)
			return
		}
		return
	}

	v := config.VOICE
	if err := v.JoinChannelAndSpeak(context.Background(), vscID, false, false); err != nil {
		logger.Err("Failed to join voice channel", err)
		return
	}
	defer v.Leave(context.Background())

}
