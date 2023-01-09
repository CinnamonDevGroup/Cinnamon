package main

import (
	"github.com/AngelFluffyOokami/Cinnamon/modules/core/commonutils"
	"github.com/AngelFluffyOokami/Cinnamon/modules/integrations/minecraft"
	minecraftdb "github.com/AngelFluffyOokami/Cinnamon/modules/integrations/minecraft/database"
)

func init() {
	for k, v := range minecraft.CommandsHandlers {
		allCommandHandlers[k] = v
	}
	allCommands = append(allCommands, minecraft.Commands...)
	DBMigrate = append(DBMigrate, MigrateDB)
	commonutils.AuthKeyUpdater = append(commonutils.AuthKeyUpdater, minecraft.RegenAuthKeys)
	for k, v := range minecraft.WebsocketHandler {
		allWebsocketHandlers[k] = v
	}
}

func MigrateDB() {
	DB.AutoMigrate(&minecraftdb.Minecraft{})
}
