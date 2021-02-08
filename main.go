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
		return
	}

	guildId := os.Getenv("GUILD_ID")

	if guildId == "" {
		log.Fatal("$GUILD_ID must be set")
		return
	}

	channelId := os.Getenv("CHANNEL_ID")

	if channelId == "" {
		log.Fatal("$CHANNEL_ID must be set")
		return
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
		return
	}

	session.AddHandler(messageCreate())

	session.AddHandler(voiceStateUpdate(channelId))

	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates)

	err = session.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	session.Close()
}

func messageCreate() func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if m.Content == "/test" {
			s.ChannelMessageSend(m.ChannelID, "シロもなかなかイケるな。しんのすけ、泣いてないで食え。")
		}

		if m.Content == "/osaka" {
			s.ChannelMessageSend(m.ChannelID, "テーマパークに来たみたいだぜ。テンション上がるなぁ～")
		}
	}
}

func voiceStateUpdate(channelId string) func(s *discordgo.Session, m *discordgo.VoiceStateUpdate) {
	return func(s *discordgo.Session, m *discordgo.VoiceStateUpdate) {
		user, err := s.User(m.UserID)
		if err != nil {
			return
		}

		before := m.BeforeUpdate.ChannelID
		after := m.VoiceState.ChannelID

		if before == "" && after != "" {
			channel, err := s.Channel(after)
			if err != nil {
				return
			}
			fmt.Println(user.Username, channel.Name)
		}

		if before != "" && after != "" {
			channel, err := s.Channel(before)
			if err != nil {
				return
			}
			fmt.Println(user.Username, channel.Name)
		}
	}
}
