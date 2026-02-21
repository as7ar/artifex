package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/as7ar/artifex/commands"
	log2 "github.com/as7ar/artifex/log"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	TOKEN string
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Something wrong when loading .env file", err)
		os.Exit(-1)
	}

	TOKEN = os.Getenv("BOT_TOKEN")

	discord, err := discordgo.New(TOKEN)

	if discord == nil {
		log.Println("Something wrong", err)
		os.Exit(-1)
	}

	discord.AddHandler(commands.PingCommand)

	_ = discord.Open()

	bot := *discord.State.User

	log.Println("Bot", log2.GREEN, bot.Username, log2.RESET, "is available")

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

}
