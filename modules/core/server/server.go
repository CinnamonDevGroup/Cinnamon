package coreserver

import (
	"github.com/AngelFluffyOokami/Cinnamon/modules/core/commonutils"
	databaseHelper "github.com/AngelFluffyOokami/Cinnamon/modules/core/database"
	coredb "github.com/AngelFluffyOokami/Cinnamon/modules/core/database/core"
	"github.com/AngelFluffyOokami/Cinnamon/modules/integrations/minecraft"
	"github.com/tjarratt/babble"

	"github.com/bwmarrin/discordgo"
)

func regenAuthKey(GID string, DB databaseHelper.DBstruct) string {

	babbler := babble.NewBabbler()

	babbler.Count = 6

	babbler.Separator = "-"

	guild := coredb.Guild{GID: GID}

	DB.Guilds.First(&guild)

	guild.AuthKey = commonutils.BabbleWords()

	DB.Guilds.Save(&guild)

	return guild.AuthKey

}

func UpdateAuthKeys(GID string, DB databaseHelper.DBstruct, AuthKey string) {

	minecraft.RegenAuthKeys(GID, DB, AuthKey)

}

func OnServerJoin(s *discordgo.Session, z *discordgo.GuildCreate, DB databaseHelper.DBstruct) {

	commonutils.CheckServerExists(z.Guild.ID, DB, s)
}

var (
	admin    int64 = discordgo.PermissionAdministrator
	Commands       = []discordgo.ApplicationCommand{
		{
			Name:                     "regenauthkey",
			Description:              "Regenerates your authentication keys for all integrations in case of being leaked.",
			DefaultMemberPermissions: &admin,
		},
		{
			Name:                     "authkey",
			Description:              "Returns the value of your auth keys.",
			DefaultMemberPermissions: &admin,
		},
	}
	CommandsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, DB databaseHelper.DBstruct){
		"regenauthkey": func(s *discordgo.Session, i *discordgo.InteractionCreate, DB databaseHelper.DBstruct) {

			commonutils.CheckServerExists(i.Interaction.ID, DB, s)
			newKeys := regenAuthKey(i.Interaction.GuildID, DB)
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Your authentication key has been regenerated, new key is as follows: \n```" + newKeys + "```\nPlease keep it safe.",
				},
			})
			if err != nil {
				panic(err)
			}

		},
		"authkey": func(s *discordgo.Session, i *discordgo.InteractionCreate, DB databaseHelper.DBstruct) {

			commonutils.CheckServerExists(i.Interaction.GuildID, DB, s)

			guild := coredb.Guild{GID: i.Interaction.GuildID}

			DB.Guilds.First(&guild)

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Your authentication key is as follows: \n```" + guild.AuthKey + "```\nPlease keep it safe.",
				},
			})
			if err != nil {
				panic(err)
			}

		},
	}
)
