package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	token := os.Getenv("TOKEN")

	if token == "" {
		log.Fatal("$TOKEN must be set")
	}

	guildId := os.Getenv("GUILD_ID")

	if guildId == "" {
		log.Fatal("GUILD_ID must be set")
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	session.AddHandler(messageCreate)

	session.Identify.Intents = discordgo.IntentsGuildMessages

	err = session.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	channels, err := session.GuildChannels(guildId)
	for _, ch := range channels {
		fmt.Println((*ch).ID)
		fmt.Println((*ch).Name)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	session.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "/test" {
		s.ChannelMessageSend(m.ChannelID, "シロもなかなかイケるな。しんのすけ、泣いてないで食え。")
	}
}
