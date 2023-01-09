package affection_db

type PerUserPreferences struct {
	UID                 string
	OverrideServerAffct []string
	AffectionAllowed    bool
}

type PerServerPreferences struct {
	AffectionAllowed bool
	GID              string
}

type AffectData struct {
	DefaultAffect     bool
	UserPreferences   []PerUserPreferences
	ServerPreferences []PerServerPreferences
}
