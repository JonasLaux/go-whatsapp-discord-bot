package main

import (
	"fmt"
	"go-whatsapp-discord-bot/whatsapp"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// global vars
var discordChannelID string
var groupJID string

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// get vars from .env
	Token := os.Getenv("DISCORD_TOKEN")
	discordChannelID = os.Getenv("DISCORD_CHANNEL_ID")
	groupJID = os.Getenv("WHATSAPP_GROUP_JID")

	print("Starting bot...\n")
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	whatsapp.Init()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Close Whatsapp
	whatsapp.Disconnect()
	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// If the message is !philipp, please throw up :)
	if m.Content == "!philipp" {

		_, err := s.ChannelMessageSend(m.ChannelID, ":face_vomiting:")
		if err != nil {
			fmt.Println(err)
		}

	}

	// If the message is in the right channel, pipe it to WhatsApp
	if m.ChannelID == discordChannelID {

		whatsapp.SendMessage(m.Author.Username+":  "+m.Content, groupJID)

	}
}
