package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/carlmjohnson/requests"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"
)

var (
	discordAuthenticationToken = os.Getenv("DISCORD_AUTHENTICATION_TOKEN")
	qaApiEndpoint              = os.Getenv("QA_API_ENDPOINT")
	qaApiKey                   = os.Getenv("QA_API_KEY")
)

type submitRequest struct {
	RequestID string `json:"request_id"`
	Query     string `json:"query"`
}

type submitResponse struct {
	RequestID string `json:"request_id"`
	Query     string `json:"query"`
	Response  string `json:"response"`
}

func handleRequest(s *discordgo.Session, m *discordgo.MessageCreate) {
	command := "!request"

	// check if command is triggered
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	if m.Author.ID == s.State.User.ID {
		return
	}

	// main
	query := strings.Replace(m.Content, command, "", 1)

	var response submitResponse
	if query != "" {
		log.Info().
			Str("author", m.Author.Username).
			Str("channel", m.Message.ChannelID).
			Msg("Query received")

		id := uuid.New()

		body := submitRequest{
			RequestID: id.String(),
			Query:     query,
		}

		err := requests.
			URL(qaApiEndpoint).
			Method(http.MethodPost).
			Path("submit").
			BodyJSON(body).
			Cookie("access_token", qaApiKey).
			ToJSON(&response).
			Fetch(context.Background())

		if err != nil {
			log.Fatal().Err(err).Msg("Failed to send request to qa-api")
		}

		log.Info().
			Str("author", m.Author.Username).
			Str("channel", m.Message.ChannelID).
			Str("response", response.Response).
			Msg(query)

		_, err = s.ChannelMessageSend(m.ChannelID, response.Response)
		if err != nil {
			log.Error().Err(err).Msg("Error sending message")
		}
	}
}

func main() {
	// init
	session, _ := discordgo.New(fmt.Sprintf("Bot %s", discordAuthenticationToken))

	// main
	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		handleRequest(s, m)
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
