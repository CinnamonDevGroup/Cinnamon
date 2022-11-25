package UserModel

type PerServerPreferences struct {
	Pronouns         string
	Nickname         string
	AffectionAllowed bool
}

type PerUserPreferences struct {
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
	ServerPrefer    []PerServerPreferences
	CurrentServers  []string
}
