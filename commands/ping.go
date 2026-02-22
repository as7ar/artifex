package commands

import (
	"fmt"
	"strconv"
	"time"

	"github.com/as7ar/noori/config"
	"github.com/as7ar/noori/embeds"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session"
)

func PingCommand(s *session.Session, c *gateway.MessageCreateEvent) {
	now := time.Now()
	messageTime := c.Timestamp.Time()
	ping := now.Sub(messageTime).Milliseconds()

	_, err := s.SendEmbedReply(
		c.ChannelID, c.Message.ID,
		embeds.New().Color(config.SymbolColor).
			Title("🏓 PONG!").Description("`"+strconv.FormatInt(ping, 10)+"ms`").
			Build())
	if err != nil {
		fmt.Println("send error:", err)
		return
	}
}
