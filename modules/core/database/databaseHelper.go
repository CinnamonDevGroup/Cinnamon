package databaseHelper

import (
	coredb "github.com/AngelFluffyOokami/Cinnamon/modules/core/database/core"
	minecraftdb "github.com/AngelFluffyOokami/Cinnamon/modules/integrations/minecraft/database"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func Init() *gorm.DB {

	cinnamondb, err := gorm.Open(sqlite.Open("database/cinnamon.db"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	cinnamondb.AutoMigrate(&coredb.Cinnamon{}, coredb.Guild{}, coredb.User{}, minecraftdb.Minecraft{})

	return cinnamondb
}
