package minecraft

import "encoding/json"

const AuthKickEvent = "playerauthkickevent"

type kickForAuth struct {
	UUID    string `json:"uuid"`
	AuthKey string `json:"authkey"`
}

const playerAuthEvent = "playerauthevent"

type playerAuthSuccessful struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
}

const notFoundEvent = "usernotfoundevent"

type kickForNotOnServer struct {
	UUID string `json:"uuid"`
}

type OutboundData struct {
	DataType string          `json:"datatype"`
	RawData  json.RawMessage `json:"rawdata"`
	API      int             `json:"version"`
}
