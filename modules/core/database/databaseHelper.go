package databaseHelper

import (
	coredb "github.com/AngelFluffyOokami/Cinnamon/modules/core/database/core"
	minecraftdb "github.com/AngelFluffyOokami/Cinnamon/modules/integrations/minecraft/database"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB DBstruct

type DBstruct struct {
	Cinnamon  *gorm.DB
	Guilds    *gorm.DB
	Users     *gorm.DB
	Minecraft *gorm.DB
}

func Init() DBstruct {

	cinnamondb, err := gorm.Open(sqlite.Open("database/cinnamon"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	guildsdb, err := gorm.Open(sqlite.Open("database/guilds"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	usersdb, err := gorm.Open(sqlite.Open("database/users"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	mcdb, err := gorm.Open(sqlite.Open("database/minecraft"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	cinnamondb.AutoMigrate(&coredb.Cinnamon{})
	guildsdb.AutoMigrate(&coredb.Guild{})
	usersdb.AutoMigrate(&coredb.User{})
	mcdb.AutoMigrate(&minecraftdb.Minecraft{})

	DB.Cinnamon = cinnamondb
	DB.Guilds = guildsdb
	DB.Users = usersdb
	DB.Minecraft = mcdb
	return DB
}
