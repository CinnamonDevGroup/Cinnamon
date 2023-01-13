package database

import (
	coredb "github.com/CinnamonDevGroup/Cinnamon/modules/core/database/core_models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {

	cinnamondb, err := gorm.Open(sqlite.Open("database/cinnamon.db"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	cinnamondb.AutoMigrate(&coredb.Cinnamon{}, coredb.Guild{}, coredb.User{}, coredb.UserModule{}, coredb.ServerModule{})

	return cinnamondb
}
