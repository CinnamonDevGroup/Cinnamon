package dbModel

import (
	guildModel "github.com/AngelFluffyOokami/Cinnamon/modules/database/GuildSide"
	userModel "github.com/AngelFluffyOokami/Cinnamon/modules/database/UserSide"
)

type Cinnamon struct {
	Guilds []guildModel.Guild
	Users  []userModel.User
}
