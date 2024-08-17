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

		// check if command is triggered
		if !strings.HasPrefix(m.Content, command) {
			return
		}

		if m.Author.ID == s.State.User.ID {
			return
		}

		// main
		question := strings.Replace(m.Content, command, "", 1)

		// send reply
		if question != "" {
			log.Info().
				Str("author", m.Author.Username).
				Str("channel", m.Message.ChannelID).
				Msg(question)

			response := question

			_, err := s.ChannelMessageSend(m.ChannelID, response)
			if err != nil {
				log.Error().Err(err).Msg("Error sending message")
			}
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
