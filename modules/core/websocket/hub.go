package websocket

import (
	"encoding/json"
	"fmt"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	Clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan []byte

	Client chan *Client

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Client:     make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

var GetWebsocketHandlers = make(chan map[string]func(receivedData IncomingData, h *Hub))
var SetWebsocketHandlers = make(chan map[string]func(receivedData IncomingData, h *Hub))

func Websocket() {
	var Handlers map[string]func(receivedData IncomingData, h *Hub)
	for {
		select {
		case Handlers = <-SetWebsocketHandlers:
		case GetWebsocketHandlers <- Handlers:
		}
	}
}

func (h *Hub) run() {
	Handlers := <-GetWebsocketHandlers
	for {
		select {
		case client := <-h.Register:

			fmt.Print("Server Client registered " + client.Addr)

			h.Clients[client] = true

		case client := <-h.Unregister:

			clientAuthKey := client.AuthKey
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}

			fmt.Print("Server Client has unregistered " + clientAuthKey)

		case data := <-h.Broadcast:
			var incomingData IncomingData
			json.Unmarshal(data, &incomingData)

			if u, ok := Handlers[incomingData.DataType]; ok {
				u(incomingData, h)
			}
		}
	}
}
