package dbModel

import (
	cinnamonModel "github.com/AngelFluffyOokami/Cinnamon/modules/database/cinnamon"
	guildModel "github.com/AngelFluffyOokami/Cinnamon/modules/database/guild"
	userModel "github.com/AngelFluffyOokami/Cinnamon/modules/database/user"
)

type Cinnamon struct {
	Guilds   []guildModel.Guild
	Users    []userModel.User
	Cinnamon cinnamonModel.Cinnamon
}
