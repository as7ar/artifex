package main

import (
	"artifex/commands"
	log2 "artifex/log"
	"log"
	"os"

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
}
