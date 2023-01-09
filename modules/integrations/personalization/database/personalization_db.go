package personalization_db

type PerServerPreferences struct {
	Pronouns Pronouns
	GID      string
	Nickname string
}

type PerUserPreferences struct {
	UID            string
	Pronouns       Pronouns
	Nickname       string
	OverrideServer []string
}

type Pronouns struct {
	Nominative string
	Objective  string
	Possessive string
}

type PersonalizationData struct {
	ServerPreferences []PerServerPreferences
	DefaultPronouns   []Pronouns
	DefaultNickname   string
}
