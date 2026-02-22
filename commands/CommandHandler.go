package commands

import (
	"strings"

	"github.com/as7ar/noori/config"
	"github.com/as7ar/noori/embeds"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func CommandHandler(c *gateway.MessageCreateEvent) {
	if c.Author.Bot {
		return
	}
	s := config.NOORI.SESSION

	if !strings.HasPrefix(c.Content, config.Prefix) {
		return
	}

	parts := strings.Fields(c.Content)
	cmd := parts[0]

	switch cmd {
	case "!ping":
		PingCommand(s, c)
	case "!radio":
		if len(parts) < 2 {
			s.SendEmbedReply(c.ChannelID, c.ID,
				embeds.New().Color(config.SymbolColor).
					Title("🚨 ERROR").Description("`!radio [<url>]`").Build())
			return
		}

		RadioCommand(s, c, parts[1])
	case "!stop":
		StopCommand(c)
	}

}
