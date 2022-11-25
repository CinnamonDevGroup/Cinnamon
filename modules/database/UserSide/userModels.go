package userModel

type perServerPreferences struct {
	Pronouns         string
	Nickname         string
	AffectionAllowed bool
}

type perUserPreferences struct {
	Pronouns            string
	Nickname            string
	AffectionAllowed    bool
	OverrideServerPron  []string
	OverrideServerNick  []string
	OverrideServerAffct []string
}

type User struct {
	DefaultPronouns string
	DefaultNickname string
	DefaultAffect   bool
	ServerPrefer    []perServerPreferences
	UserPrefer      []perUserPreferences
	CurrentServers  []string
	XP              []globalXP
}

type globalXP struct {
	TotalXP     int
	PerServerXP []guildXP
}

type guildXP struct {
	GID string
	XP  int
}
