package guildModel

import MinecraftModel "github.com/AngelFluffyOokami/Cinnamon/modules/database/minecraft"

type config struct {
	Minecraft           MinecraftModel.Minecraft
	DefaultAdminChannel string
}

type Guild struct {
	GID    string
	Joined string
	About  information `gorm:"serializer:json"`
	Config config      `gorm:"serializer:json"`
}

type information struct {
	JoinedAt     []string
	LeftAt       []string
	UserAmount   int
	UserInfo     []guildUser
	MessageCount []message
}

type message struct {
	MessageCount int
	TimeCount    int
}

type guildUser struct {
	UID             string
	XP              int
	JoinedPositions []int
	LeftPositions   []int
	Moderated       []userModeration
}

type userModeration struct {
	Warnings   []warning
	Mutes      []mute
	Kicks      []kick
	Bans       []ban
	ActiveBan  bool
	ActiveMute bool
}

type warning struct {
	WarningReason string
	WarnedAt      int
}

type ban struct {
	BanReason string
	BannedOn  int
}

type kick struct {
	KickReason string
	KickedOn   int
}

type mute struct {
	MuteReason    string
	MuteDuration  int
	MuteStartedAt int
	MuteExpired   bool
}
