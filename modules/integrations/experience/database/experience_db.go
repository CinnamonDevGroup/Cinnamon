package experience_db

type XPData struct {
	TotalXP     int
	PerServerXP []GuildXP
}

type GuildXP struct {
	GID string
	XP  int
}

type TotalXP struct {
	TotalXP int64
}
