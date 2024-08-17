package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"
)

var (
	discordAuthenticationToken = os.Getenv("DISCORD_AUTHENTICATION_TOKEN")
)

func main() {
	// init
	session, _ := discordgo.New(fmt.Sprintf("Bot %s", discordAuthenticationToken))

	// main
	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		command := "!request"
		args := strings.Split(m.Content, " ")
		if args[0] != command {
			return
		}

		if m.Author.ID == s.State.User.ID {
			return
		}

		input := m.Content[1:]
		response := input

		_, err := s.ChannelMessageSend(m.ChannelID, response)
		if err != nil {
			log.Fatal().Err(err).Msg("Error sending message")
		}
	})

	// bot init
	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Info().Msgf("Logged in as %s", r.User.String())
	})

	// start bot
	err := session.Open()
	if err != nil {
		log.Fatal().Err(err).Msg("Could not open session")
	}

	// keep the bot running
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch

	err = session.Close()
	if err != nil {
		log.Error().Err(err).Msg("Could not close session gracefully")
	}
}
