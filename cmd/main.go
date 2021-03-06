package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/JoungSik/gambler/cmd/models"
	"github.com/JoungSik/gambler/configs"
	"github.com/bwmarrin/discordgo"
)

func main() {
	config := configs.NewConfig()

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + config.TOKEN)
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

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

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

	if m.Content == "" {
		return
	}

	if m.Content[0] == '!' {
		db, err := configs.InitDB(false)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}

		database, _ := db.DB()
		defer database.Close()

		command, err := models.Parse(m.Content, m)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}

		s.ChannelMessageSend(m.ChannelID, command.Execute(db))
	}
}
