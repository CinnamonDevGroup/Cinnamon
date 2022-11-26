package minecraft

import (
	"encoding/json"
	"fmt"

	databaseHelper "github.com/AngelFluffyOokami/Cinnamon/modules/core/database"
	"github.com/bwmarrin/discordgo"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run(s *discordgo.Session, DB databaseHelper.DBstruct) {
	for {
		select {
		case client := <-h.register:
			clientAuthId := client.AuthID

			fmt.Print("Server Client registered " + clientAuthId)

			h.clients[client] = true

		case client := <-h.unregister:

			clientAuthId := client.AuthID
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

			fmt.Print("Server Client has unregistered " + clientAuthId)

		case data := <-h.broadcast:
			fmt.Print("msg received")

			var fauxMessage discordMessage

			fauxMessage.Channel = "exchannel"
			fauxMessage.Mention = "exmention"
			fauxMessage.Message = "exmessage"
			fauxMessage.User = "exuser"
			fauxjson, err := json.Marshal(fauxMessage)

			if err != nil {
				panic(err)
			}

			var receivedData IncomingData
			json.Unmarshal(data, &receivedData)
			authID := receivedData.AuthID
			for client := range h.clients {

				if authID == client.AuthID {

					switch receivedData.DataType {
					case "playerMessage":
						onPlayerMessage(receivedData.Data, s)
					case "playerJoin":
						onPlayerJoin(receivedData.Data, s)
					}
					select {
					case client.send <- fauxjson:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}

			}
		}
	}
}
