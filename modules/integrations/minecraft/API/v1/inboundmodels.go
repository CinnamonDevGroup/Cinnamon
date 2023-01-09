package minecraft_api_v1

import (
	"encoding/json"
)

type ChatMessage struct {
	UUID    string `json:"uuid"`
	Message string `json:"message"`
	Mention string `json:"mention"`
	Channel string `json:"channel"`
}

type PlayerJoin struct {
	UUID string `json:"uuid"`
}

type Data struct {
	DataType string          `json:"datatype"`
	RawData  json.RawMessage `json:"rawdata"`
	AuthKey  string          `json:"authkey"`
	GID      string          `json:"gid"`
}

type Authenticate struct {
	AuthKey        string `json:"authkey"`
	DefaultChannel string `json:"channel"`
	GuildID        string `json:"guild"`
}

type ConnectionStatus struct {
	AuthKey string `json:"authkey"`
	GID     string `json:"gid"`
	Status  int    `json:"status"`
}

type DiscordMessage struct {
	User    string `json:"user"`
	Mention string `json:"mention"`
	Message string `json:"message"`
	Channel string `json:"channel"`
}
