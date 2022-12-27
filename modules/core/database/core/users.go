package coredb

import "encoding/json"

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
	UID             string                 `gorm:"primaryKey"`
	ServerPrefer    []perServerPreferences `gorm:"serializer:json"`
	UserPrefer      []perUserPreferences   `gorm:"serializer:json"`
	CurrentServers  []string               `gorm:"serializer:json"`
	XP              []globalXP             `gorm:"serializer:json"`
	Services        []Service              `gorm:"serializer:json"`
}

type Service struct {
	Service string `gorm:"primaryKey"`
	UUID    string `gorm:"primaryKey"`
	Data    json.RawMessage
	UID     string
	AuthKey string
}

type globalXP struct {
	TotalXP     int
	PerServerXP []guildXP
}

type guildXP struct {
	GID string
	XP  int
}
