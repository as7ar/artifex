package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/as7ar/noori/commands"
	"github.com/as7ar/noori/config"
	logger2 "github.com/as7ar/noori/logger"
	noori2 "github.com/as7ar/noori/noori"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/state/store/defaultstore"
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

	st := session.NewWithIntents("Bot "+config.TOKEN,
		gateway.IntentGuilds, gateway.IntentGuildVoiceStates, gateway.IntentGuildMessages)

	voice.AddIntents(st)

	st.AddHandler(commands.CommandHandler)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := st.Open(ctx); err != nil {
		log.Fatalln("Failed to connect:", err)
	}
	defer st.Close()

	me, _ := st.Me()
	logger2.Info("Bot ", me.Username, " is available")

	config.NOORI = &noori2.App{
		SESSION: st,
		STATE:   state.NewFromSession(st, defaultstore.New()),
	}

	<-ctx.Done()
}
