package websocket

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub *Hub

	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte

	Authenticated bool

	User
}

const APINegotiate = "negotiateapi"

const maxAPI = 1

type NegotiateAPI struct {
	APIVersion int `json:"version"`
}

type ConnectionStatus struct {
	AuthKey string `json:"authkey"`
	GID     string `json:"gid"`
	Status  int    `json:"status"`
}

type IncomingData struct {
	DataType   string          `json:"datatype"`
	RawData    json.RawMessage `json:"rawdata"`
	APIVersion int             `json:"version"`
	Client     *Client
}
type OutboundData struct {
	DataType   string          `json:"datatype"`
	RawData    json.RawMessage `json:"rawdata"`
	APIVersion int             `json:"version"`
	UUID       string
}
type User struct {
	AuthKey    string
	APIVersion int
	Service    string
	Addr       string
	EnterAt    time.Time
	UUID       string
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
)

var (
	newline = []byte{'\n'}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
