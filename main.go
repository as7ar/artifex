package main

import (
	"context"
	"log"
	"os"

	"github.com/as7ar/noori/commands"
	"github.com/as7ar/noori/config"
	logger2 "github.com/as7ar/noori/logger"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session"
	"github.com/diamondburned/arikawa/v3/voice"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		logger2.Err("Something wrong when loading .env file", err)
		os.Exit(-1)
	}

	config.TOKEN = os.Getenv("BOT_TOKEN")

	s := session.NewWithIntents("Bot "+config.TOKEN,
		gateway.IntentGuilds, gateway.IntentGuildVoiceStates, gateway.IntentGuildMessages)

	if s == nil {
		logger2.Err("Something wrong", err)
		os.Exit(-1)
	}

	s.AddHandler(func(c *gateway.MessageCreateEvent) {
		if c.Author.Bot {
			return
		}
		logger2.Debug("received:", c.Content)
	})
	s.AddHandler(commands.CommandHandler)

	if err := s.Open(context.Background()); err != nil {
		log.Fatalln("Failed to connect:", err)
	}
	defer func(s *session.Session) {
		err := s.Close()
		if err != nil {
			logger2.Err(err)
		}
	}(s)

	u, err := s.Me()
	if err != nil {
		log.Fatalln("Failed to get myself:", err)
	}

	if u != nil {
		log.Println("Bot", u.Username, "is available")
	}

	v, err := voice.NewSession(s)
	if err != nil {
		log.Fatalln("failed to create voice session:", err)
	}

	config.VOICE = v
	config.DISCORD = s

	select {}
}
