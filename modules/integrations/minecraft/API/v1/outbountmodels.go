package minecraft_api_v1

import "encoding/json"

const AuthKickEvent = "playerauthkickevent"

type KickForAuth struct {
	UUID    string `json:"uuid"`
	AuthKey string `json:"authkey"`
}

const PlayerAuthEvent = "playerauthevent"

type PlayerAuthSuccessful struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
}

const NotFoundEvent = "usernotfoundevent"

type KickForNotOnServer struct {
	UUID string `json:"uuid"`
}

type OutboundData struct {
	DataType string          `json:"datatype"`
	RawData  json.RawMessage `json:"rawdata"`
	API      int             `json:"version"`
}
