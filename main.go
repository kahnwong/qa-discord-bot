package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

var (
	discordAuthenticationToken = os.Getenv("DISCORD_AUTHENTICATION_TOKEN")
)

func main() {
	// init
	session, _ := discordgo.New(fmt.Sprintf("Bot %s", discordAuthenticationToken))

	// main
	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		fmt.Println(m.Content)

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
			log.Fatal(err)
		}
	})

	// bot init
	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as %s", r.User.String())
	})

	// start bot
	err := session.Open()
	if err != nil {
		log.Fatalf("could not open session: %s", err)
	}

	// keep the bot running
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch

	err = session.Close()
	if err != nil {
		log.Printf("could not close session gracefully: %s", err)
	}
}
