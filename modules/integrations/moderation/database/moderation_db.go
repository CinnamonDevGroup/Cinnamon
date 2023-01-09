package moderation_db

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
