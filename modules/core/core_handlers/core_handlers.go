package core_handlers

import (
	"github.com/CinnamonDevGroup/Cinnamon/modules/core/common"
	"github.com/CinnamonDevGroup/Cinnamon/modules/core/database/core_models"
	"github.com/bwmarrin/discordgo"
)

func regenAuthKey(GID string) string {

	DB := common.DB
	guild := core_models.Guild{GID: GID}

	DB.First(&guild)

	oldKey := guild.AuthKey

	guild.AuthKey = common.BabbleWords()

	DB.Save(&guild)

	updateDBAuthKeys(GID, guild.AuthKey, oldKey)

	return guild.AuthKey

}

func updateDBAuthKeys(GID string, AuthKey string, OldKey string) {

	//TODO AuthKeyUpdater

	for _, x := range common.AuthKeyUpdater {
		x(GID, AuthKey, OldKey)
	}

}

func OnServerJoin(z *discordgo.GuildCreate) {

	common.CheckGuildExists(z.Guild.ID)
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
			s := common.Session
			common.CheckGuildExists(i.Interaction.GuildID)
			newKeys := regenAuthKey(i.Interaction.GuildID)
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Your authentication key has been regenerated, new key is as follows: \n```" + newKeys + "```\nPlease keep it safe.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				panic(err)
			}

		},
		"authkey": func(i *discordgo.InteractionCreate) {
			s := common.Session
			DB := common.DB

			common.CheckGuildExists(i.Interaction.GuildID)

			guild := core_models.Guild{GID: i.Interaction.GuildID}

			DB.First(&guild)

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Your authentication key is as follows: \n```" + guild.AuthKey + "```\nPlease keep it safe.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				panic(err)
			}

		},
	}
)
