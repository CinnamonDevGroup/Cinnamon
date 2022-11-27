package minecraft

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	client chan *Client

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		client:     make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run(s *discordgo.Session, DB *gorm.DB) {
	for {
		select {
		case client := <-h.register:

			fmt.Print("Server Client registered " + client.Addr)

			h.clients[client] = true

		case client := <-h.unregister:

			clientAuthKey := client.AuthKey
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

			fmt.Print("Server Client has unregistered " + clientAuthKey)

		case data := <-h.broadcast:
			fmt.Print("msg received")
			var receivedData IncomingData
			json.Unmarshal(data, &receivedData)
			client := <-h.client

			if client.Authenticated {

				switch receivedData.DataType {
				case "playerMessage":
					onPlayerMessage(receivedData.Data, s)
				case "playerJoin":
					onPlayerJoin(receivedData.Data, s)

				}

			} else {
				switch receivedData.DataType {
				case "authenticate":
					clientAuthenticate(client, DB, s, h, receivedData.Data)
				default:
					connection := ConnectionStatus{
						AuthKey: "",
						GID:     "",
						Status:  http.StatusUnauthorized,
					}
					response, err := json.Marshal(connection)
					if err != nil {
						panic(err)
					}
					client.send <- response

					close(client.send)
					delete(h.clients, client)

				}
			}

		}
	}
}
