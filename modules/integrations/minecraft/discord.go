package minecraft

import (
	"errors"

	"github.com/AngelFluffyOokami/Cinnamon/modules/core/commonutils"

	coredb "github.com/AngelFluffyOokami/Cinnamon/modules/core/database/core"
	minecraftdb "github.com/AngelFluffyOokami/Cinnamon/modules/integrations/minecraft/database"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

func checkGuildExists(GID string) bool {
	DB := <-commonutils.GetDB

	guild := coredb.Guild{GID: GID}
	DB.First(&guild)
	server := minecraftdb.Minecraft{GID: GID, AuthKey: guild.AuthKey}

	result := DB.First(&server)

	notFound := errors.Is(result.Error, gorm.ErrRecordNotFound)

	var guildExists bool
	if notFound {
		guildExists = false
	} else {
		guildExists = true
	}

	return guildExists

}

func initializeGuild(GID string) {

	DB := <-commonutils.GetDB

	guild := coredb.Guild{GID: GID}

	DB.First(&guild)

	server := minecraftdb.Minecraft{
		GID:     GID,
		AuthKey: guild.AuthKey,
		Active:  true,
	}

	DB.Create(&server)

}

func unlinkServer(GID string) {

	DB := <-commonutils.GetDB
	guild := coredb.Guild{GID: GID}

	DB.First(&guild)
	server := minecraftdb.Minecraft{GID: GID, AuthKey: guild.AuthKey}

	DB.First(&server)

	server.Active = false

	DB.Save(&server)
}

func RegenAuthKeys(GID string, AuthKey string, OldKey string) {

	DB := <-commonutils.GetDB

	guild := coredb.Guild{GID: GID}
	DB.First(&guild)
	server := minecraftdb.Minecraft{GID: GID, AuthKey: OldKey}
	result := DB.First(&server)

	notFound := errors.Is(result.Error, gorm.ErrRecordNotFound)

	var guildExists bool
	if notFound {
		guildExists = false
	} else {
		guildExists = true
	}

	if guildExists {
		server := minecraftdb.Minecraft{GID: GID}

		DB.First(&server)

		oldServer := server
		server.AuthKey = AuthKey

		DB.Save(&server)

		DB.Delete(oldServer)

	} else {
		return
	}

}

func deleteGuildData(GID string) {

	DB := <-commonutils.GetDB
	guild := coredb.Guild{GID: GID}

	DB.First(&guild)
	server := minecraftdb.Minecraft{
		GID:     GID,
		AuthKey: guild.AuthKey,
	}

	DB.Delete(&server)
}

func enableGuild(GID string) string {
	DB := <-commonutils.GetDB
	guild := coredb.Guild{GID: GID}

	DB.First(&guild)

	server := minecraftdb.Minecraft{GID: GID, AuthKey: guild.AuthKey}

	DB.First(&server)
	var message string
	if !server.Active {
		server.Active = true
		message = "Minecraft integration enabled.\n"
	} else {
		message = "Minecraft integration already enabled.\n"
	}

	DB.Save(&server)

	return message
}

var (
	admin    int64 = discordgo.PermissionAdministrator
	Commands       = []discordgo.ApplicationCommand{
		{
			Name:                     "linkminecraftserver",
			Description:              "Link and enable Minecraft integration functionality to the bot.",
			DefaultMemberPermissions: &admin,
		},
		{
			Name:                     "unlinkminecraftserver",
			Description:              "Disables the Minecraft integration related functionality, keeping player data intact.",
			DefaultMemberPermissions: &admin,
		},
		{
			Name:        "deleteminecraftlink",
			Description: "Deletes all Minecraft integration data from database and disables Minecraft related functionality.",
		},
	}
	CommandsHandlers = map[string]func(i *discordgo.InteractionCreate){
		"linkminecraftserver": func(i *discordgo.InteractionCreate) {
			DB := <-commonutils.GetDB
			s := <-commonutils.GetSession
			commonutils.CheckGuildExists(i.Interaction.GuildID)

			exists := checkGuildExists(i.Interaction.GuildID)

			var message string
			if !exists {
				initializeGuild(i.Interaction.GuildID)
				message = "Minecraft integration enabled.\n"
			} else {
				message = enableGuild(i.Interaction.GuildID)
			}

			guild := coredb.Guild{GID: i.Interaction.GuildID}

			DB.First(&guild)

			server := minecraftdb.Minecraft{GID: i.Interaction.GuildID, AuthKey: guild.AuthKey}

			DB.First(&server)

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: message + "To link your Minecraft server to this guild, please add the following authentication key to the configuration file created by the mod.\nYour secret server Authentication Key is: \n```" + guild.AuthKey + "```\nPlease keep this key safe.\nMinecraft functionality has been enabled for this guild.",
					Flags:   1 << 6,
				},
			})
			if err != nil {
				panic(err)
			}

		},
		"deleteminecraftlink": func(i *discordgo.InteractionCreate) {
			s := <-commonutils.GetSession
			commonutils.CheckGuildExists(i.Interaction.GuildID)
			exists := checkGuildExists(i.Interaction.GuildID)

			if !exists {
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "No Minecraft server has been linked to this guild yet.",
						Flags:   1 << 6,
					},
				})
				if err != nil {
					panic(err)
				}
			} else {
				deleteGuildData(i.Interaction.GuildID)

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Minecraft integration data has been deleted.",
						Flags:   1 << 6,
					},
				})
				if err != nil {
					panic(err)
				}

			}

		},

		"unlinkminecraftserver": func(i *discordgo.InteractionCreate) {
			s := <-commonutils.GetSession
			commonutils.CheckGuildExists(i.Interaction.GuildID)
			exists := checkGuildExists(i.Interaction.GuildID)

			if exists {

				exists = checkGuildExists(i.Interaction.GuildID)
				if exists {
					unlinkServer(i.Interaction.GuildID)

					err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Minecraft server has been unlinked from this guild.",
							Flags:   1 << 6,
						},
					})
					if err != nil {
						panic(err)
					}
				} else {

					err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Minecraft server has already been unlinked.",
							Flags:   1 << 6,
						},
					})
					if err != nil {
						panic(err)
					}

				}
			} else {

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "No Minecraft server has been linked to this guild yet.",
						Flags:   1 << 6,
					},
				})
				if err != nil {
					panic(err)
				}

			}
		},
	}
)
