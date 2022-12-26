package coreserver

import (
	"github.com/AngelFluffyOokami/Cinnamon/modules/core/commonutils"
	coredb "github.com/AngelFluffyOokami/Cinnamon/modules/core/database/core"

	"github.com/bwmarrin/discordgo"
)

func regenAuthKey(GID string) string {

	DB := <-commonutils.GetDB
	guild := coredb.Guild{GID: GID}

	DB.First(&guild)

	oldKey := guild.AuthKey

	guild.AuthKey = commonutils.BabbleWords()

	DB.Save(&guild)

	UpdateAuthKeys(GID, guild.AuthKey, oldKey)

	return guild.AuthKey

}

func UpdateAuthKeys(GID string, AuthKey string, OldKey string) {

	//TODO AuthKeyUpdater

	for _, x := range commonutils.AuthKeyUpdater {
		x(GID, AuthKey, OldKey)
	}

}

func OnServerJoin(z *discordgo.GuildCreate) {

	commonutils.CheckGuildExists(z.Guild.ID)
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
	CommandsHandlers = map[string]func(i *discordgo.InteractionCreate){
		"regenauthkey": func(i *discordgo.InteractionCreate) {
			s := <-commonutils.GetSession
			commonutils.CheckGuildExists(i.Interaction.GuildID)
			newKeys := regenAuthKey(i.Interaction.GuildID)
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
		"authkey": func(i *discordgo.InteractionCreate) {
			s := <-commonutils.GetSession
			DB := <-commonutils.GetDB

			commonutils.CheckGuildExists(i.Interaction.GuildID)

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
