package minecraft

import (
	"errors"

	"github.com/AngelFluffyOokami/Cinnamon/modules/core/commonutils"
	databaseHelper "github.com/AngelFluffyOokami/Cinnamon/modules/core/database"
	coredb "github.com/AngelFluffyOokami/Cinnamon/modules/core/database/core"
	minecraftdb "github.com/AngelFluffyOokami/Cinnamon/modules/integrations/minecraft/database"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

func checkExists(GID string, DB databaseHelper.DBstruct) bool {

	var guild = minecraftdb.Minecraft{GID: GID}

	result := DB.Minecraft.First(&guild)

	Found := errors.Is(result.Error, gorm.ErrRecordNotFound)

	return Found

}

func initializeServer(GID string, DB databaseHelper.DBstruct) {

	guild := coredb.Guild{GID: GID}

	DB.Guilds.First(&guild)

	server := minecraftdb.Minecraft{
		GID:     GID,
		AuthKey: guild.AuthKey,
	}

	DB.Minecraft.Create(&server)

}

func unlinkServer(GID string, DB databaseHelper.DBstruct) {
	server := minecraftdb.Minecraft{GID: GID}

	DB.Minecraft.Delete(&server)
}

func RegenAuthKeys(GID string, DB databaseHelper.DBstruct, AuthKey string) {

	exists := checkExists(GID, DB)

	if exists {
		server := minecraftdb.Minecraft{GID: GID}

		DB.Minecraft.First(&server)

		server.AuthKey = AuthKey

		DB.Minecraft.Save(&server)
	} else {
		return
	}

}

var (
	admin    int64 = discordgo.PermissionAdministrator
	Commands       = []discordgo.ApplicationCommand{
		{
			Name:                     "linkminecraftserver",
			Description:              "Link your Minecraft server to the bot.",
			DefaultMemberPermissions: &admin,
		},
		{
			Name:                     "unlinkminecraftserver",
			Description:              "Unlinks your Minecraft server from this guild, disabling all Minecraft related functionality",
			DefaultMemberPermissions: &admin,
		},
	}
	CommandsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, DB databaseHelper.DBstruct){
		"linkminecraftserver": func(s *discordgo.Session, i *discordgo.InteractionCreate, DB databaseHelper.DBstruct) {
			commonutils.CheckServerExists(i.Interaction.ID, DB, s)
			notexists := checkExists(i.Interaction.GuildID, DB)

			if notexists {
				initializeServer(i.Interaction.GuildID, DB)
			}

			var guild = minecraftdb.Minecraft{GID: i.Interaction.GuildID}

			DB.Minecraft.First(&guild)
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "To link your Minecraft server to this guild, please add the following authentication key to the configuration file created by the mod.\nYour secret server Authentication Key is: \n```" + guild.AuthKey + "```\nPlease keep this key safe.\nMinecraft functionality has been enabled for this guild.",
				},
			})
			if err != nil {
				panic(err)
			}

		},

		"unlinkminecraftserver": func(s *discordgo.Session, i *discordgo.InteractionCreate, DB databaseHelper.DBstruct) {
			commonutils.CheckServerExists(i.Interaction.ID, DB, s)
			notexists := checkExists(i.Interaction.GuildID, DB)

			if notexists {

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "No Minecraft server has been linked to this guild yet.",
					},
				})
				if err != nil {
					panic(err)
				}
			} else {

				unlinkServer(i.Interaction.GuildID, DB)

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Minecraft server has been unlinked from this guild.",
					},
				})
				if err != nil {
					panic(err)
				}

			}
		},
	}
)
