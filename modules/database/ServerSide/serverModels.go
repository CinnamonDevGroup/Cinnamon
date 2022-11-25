package databaseModels

type Config struct {
	Minecraft           Minecraft
	DefaultAdminChannel string
}

type Server struct {
	GID    string
	Joined string
	About  Information
	Config Config
}

type Information struct {
	JoinedAt     []string
	LeftAt       []string
	UserAmount   int
	UserInfo     []ServerUser
	MessageCount []Message
}

type Message struct {
	MessageTime string
	UID         string
}

type ServerUser struct {
	UID             string
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
	ActiveBan  bool
	ActiveMute bool
}

type Warning struct {
	WarningReason string
	WarnedAt      int
}

type Ban struct {
	BanReason string
	BannedAt  int
}

type Kick struct {
	KickReason string
	KickedOn   int
}
type Mute struct {
	MuteReason    string
	MuteDuration  int
	MuteStartedAt int
	MuteExpired   bool
}
