//go:build minecraft
// +build minecraft

package main

import (
	"github.com/CinnamonDevGroup/Cinnamon/modules/core/common"
	minecraft_db "github.com/CinnamonDevGroup/Cinnamon/modules/integrations/minecraft/database"
	minecraft_discord "github.com/CinnamonDevGroup/Cinnamon/modules/integrations/minecraft/discord"
	minecraft_websocket "github.com/CinnamonDevGroup/Cinnamon/modules/integrations/minecraft/websocket"
)

func init() {
	for k, v := range minecraft_discord.CommandsHandlers {
		allCommandHandlers[k] = v
	}
	allCommands = append(allCommands, minecraft_discord.Commands...)
	DBMigrate = append(DBMigrate, MigrateDB)
	common.AuthKeyUpdater = append(common.AuthKeyUpdater, minecraft_discord.RegenAuthKeys)
	for k, v := range minecraft_websocket.WebsocketHandler {
		allWebsocketHandlers[k] = v
	}
}

func MigrateDB() {
	DB.AutoMigrate(&minecraft_db.Minecraft{})
}
