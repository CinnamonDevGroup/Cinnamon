package coreserver

import (
	"github.com/AngelFluffyOokami/Cinnamon/modules/core/commonutils"
	coredb "github.com/AngelFluffyOokami/Cinnamon/modules/core/database/core"
	"github.com/AngelFluffyOokami/Cinnamon/modules/integrations/minecraft"
	"gorm.io/gorm"

	"github.com/bwmarrin/discordgo"
)

func regenAuthKey(GID string, DB *gorm.DB) string {

	guild := coredb.Guild{GID: GID}

	DB.First(&guild)

	oldKey := guild.AuthKey

	guild.AuthKey = commonutils.BabbleWords()

	DB.Save(&guild)

	UpdateAuthKeys(GID, DB, guild.AuthKey, oldKey)

	return guild.AuthKey

}

func UpdateAuthKeys(GID string, DB *gorm.DB, AuthKey string, OldKey string) {

	minecraft.RegenAuthKeys(GID, DB, AuthKey, OldKey)

}

func OnServerJoin(s *discordgo.Session, z *discordgo.GuildCreate, DB *gorm.DB) {

	commonutils.CheckGuildExists(z.Guild.ID, DB, s)
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
	CommandsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, DB *gorm.DB){
		"regenauthkey": func(s *discordgo.Session, i *discordgo.InteractionCreate, DB *gorm.DB) {

			commonutils.CheckGuildExists(i.Interaction.GuildID, DB, s)
			newKeys := regenAuthKey(i.Interaction.GuildID, DB)
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Your authentication key has been regenerated, new key is as follows: \n```" + newKeys + "```\nPlease keep it safe.",
					Flags:   1 << 6,
				},
			})
			if err != nil {
				panic(err)
			}

		},
		"authkey": func(s *discordgo.Session, i *discordgo.InteractionCreate, DB *gorm.DB) {

			commonutils.CheckGuildExists(i.Interaction.GuildID, DB, s)

			guild := coredb.Guild{GID: i.Interaction.GuildID}

			DB.First(&guild)

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Your authentication key is as follows: \n```" + guild.AuthKey + "```\nPlease keep it safe.",
					Flags:   1 << 6,
				},
			})
			if err != nil {
				panic(err)
			}

		},
	}
)
