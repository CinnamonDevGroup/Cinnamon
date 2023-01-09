package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CinnamonDevGroup/Cinnamon/modules/core/common"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	Clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan IncomingData

	Client chan *Client

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client
}

var WriteToWebsocket = make(chan ClientCache)

func newHub() *Hub {
	return &Hub{
		Broadcast:  make(chan IncomingData),
		Register:   make(chan *Client),
		Client:     make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

var WebsocketHandlers map[string]func(receivedData IncomingData, h *Hub)

type ClientCache struct {
	Client         *Client
	UUID           string
	OutboundData   OutboundData
	GID            string
	DefaultChannel string
	AuthKey        string
	Service        string
}

var ClientsCache []ClientCache

func (h *Hub) run() {
	Handlers := WebsocketHandlers
	for {
		select {
		case client := <-h.Register:

			fmt.Print("Server Client registered " + client.Addr)

			NewClient := ClientCache{
				Client: client,
				UUID:   client.User.UUID,
			}

			ClientsCache = append(ClientsCache, NewClient)

			h.Clients[client] = true

		case client := <-h.Unregister:

			clientAuthKey := client.AuthKey
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}

			fmt.Print("Server Client has unregistered " + clientAuthKey)

		case data := <-h.Broadcast:
			if data.APIVersion == 0 {

				if data.DataType == APINegotiate {
					var apiResponse NegotiateAPI
					var api NegotiateAPI
					json.Unmarshal(data.RawData, &api)
					if maxAPI <= api.APIVersion {
						apiResponse.APIVersion = maxAPI

					} else if maxAPI > api.APIVersion {
						apiResponse.APIVersion = api.APIVersion

					}
					response, _ := json.Marshal(apiResponse)
					data.Client.Send <- response
				} else {
					connection := ConnectionStatus{
						AuthKey: "",
						GID:     "",
						Status:  http.StatusBadRequest,
					}
					response, _ := json.Marshal(connection)
					data.Client.Send <- response

					close(data.Client.Send)
					delete(h.Clients, data.Client)
				}
			}

			if u, ok := Handlers[data.DataType]; ok {
				u(data, h)
			}
		case outdata := <-WriteToWebsocket:
			jsonvar, err := json.Marshal(outdata.OutboundData)
			if err != nil {
				common.LogEvent("Websocket Client Send Error Event: "+fmt.Sprint(err), common.LogError)
			}
			outdata.Client.Send <- jsonvar
		}
	}
}
