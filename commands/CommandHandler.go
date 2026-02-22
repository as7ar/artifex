package commands

import (
	"strings"

	"github.com/as7ar/noori/config"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func CommandHandler(c *gateway.MessageCreateEvent) {
	if c.Author.Bot {
		return
	}
	s := config.DISCORD

	if !strings.HasPrefix(c.Content, config.Prefix) {
		return
	}

	switch c.Content {
	case "!ping":
		PingCommand(s, c)
	case "!radio":
		RadioCommand(s, c)
	}
}
