package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	voice *discordgo.VoiceConnection
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	Token := os.Getenv("TOKEN")

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsAll
	

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
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

	perms, err := s.State.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if err != nil {
		fmt.Println("Error getting user channel permissions,", err)
		return
	}

	if (perms & discordgo.PermissionAdministrator) > 0 {
		// if m.Content == "!join" {
		// 	voiceState, err := s.State.VoiceState(m.GuildID, m.Author.ID)
		// 	if err != nil {
		// 		return
		// 	}

		// 	voiceChannel, err := s.State.Channel(voiceState.ChannelID)
		// 	if err != nil {
		// 		return
		// 	}

		// 	// Print the user's voice channel name
		// 	fmt.Printf("User %s is in voice channel %s\n", m.Author.Username, voiceChannel.Name)

		// 	voice, err = s.ChannelVoiceJoin(m.GuildID, voiceChannel.ID, false, false)
		// 	if err != nil {
		// 		fmt.Println("failed to join voice channel:", err)
		// 		return
		// 	}
		// }

		// if m.Content == "!leave" {
		// 	_, err := s.ChannelVoiceJoin(m.GuildID, "", false, true)
		// 	if err != nil {
		// 		return 
		// 	}
		// }

		if m.Content == "!m" {
			mute(s, m, true)
		}

		if m.Content == "!u" {
			mute(s, m, false)
		}

		if strings.HasPrefix(m.Content, "!") {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
		}
		
	}
}

func mute(s *discordgo.Session, m *discordgo.MessageCreate, mute bool) {
	userVoiceState, err := s.State.VoiceState(m.GuildID, m.Author.ID)
	if err != nil {
		return
	}

	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		return 
	}

	for _, voiceState := range guild.VoiceStates {
		if voiceState.ChannelID == userVoiceState.ChannelID {
			s.GuildMemberMute(m.GuildID, voiceState.UserID, mute)
		}
	}
}