package discord_client

import (
	"log"

	"github.com/CinnamonDevGroup/Cinnamon/modules/core/common"
	"github.com/bwmarrin/discordgo"
)

func Init(disToken string) *discordgo.Session {
	s, err := discordgo.New("Bot " + disToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err = s.Open()
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	common.Session = s
	return s

}
