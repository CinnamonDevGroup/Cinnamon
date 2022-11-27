package minecraft

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

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

type User struct {
	AuthKey string
	Addr    string
	EnterAt time.Time
}

type chatMessage struct {
	Player  string `json:"player"`
	Message string `json:"message"`
	Mention string `json:"mention"`
	Channel string `json:"channel"`
	Guild   string `json:"guild"`
}

type playerJoin struct {
	Username string `json:"player"`
	GID      string `json:"gid"`
	UUID     string `json:"uuid"`
}

type Multidbmodel struct {
	User     *gorm.DB
	Guild    *gorm.DB
	Cinnamon *gorm.DB
}

type playerAuthKeySend struct {
	AuthKey string `json:"authkey"`
}

type IncomingData struct {
	DataType string          `json:"datatype"`
	Data     json.RawMessage `json:"data"`
	AuthKey  string          `json:"authkey"`
}

type authenticate struct {
	AuthKey        string `json:"authkey"`
	DefaultChannel string `json:"channel"`
	GuildID        string `json:"guild"`
}

type ConnectionStatus struct {
	AuthKey string `json:"authkey"`
	GID     string `json:"gid"`
	Status  int    `json:"status"`
}

type discordMessage struct {
	User    string `json:"user"`
	Mention string `json:"mention"`
	Message string `json:"message"`
	Channel string `json:"channel"`
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	Authenticated bool

	User
}
