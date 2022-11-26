package coredb

type Config struct {
	DefaultAdminChannel string
}

type Guild struct {
	GID     string `gorm:"primaryKey"`
	Joined  int64
	AuthKey string
	About   Information `gorm:"serializer:json"`
	Config  Config      `gorm:"serializer:json"`
}

type Information struct {
	JoinedAt     []int64
	LeftAt       []int64
	UserAmount   int
	UserInfo     []GuildUser
	MessageCount []Message
}

type Message struct {
	MessageCount int
	TimeCount    int64
}

type GuildUser struct {
	UID             string
	UUID            string
	XP              int
	JoinedPositions []int
	LeftPositions   []int
	Moderated       []UserModeration
}

type UserModeration struct {
	Warnings   []Warning
	Mutes      []Mute
	Kicks      []Kick
	Bans       []Ban
	IsMember   bool
	ActiveBan  bool
	ActiveMute bool
}

type Warning struct {
	WarningReason string
	WarnedAt      int64
}

type Ban struct {
	BanReason string
	BannedOn  int64
}

type Kick struct {
	KickReason string
	KickedOn   int64
}

type Mute struct {
	MuteReason    string
	MuteDuration  int64
	MuteStartedAt int64
	MuteExpired   bool
}
