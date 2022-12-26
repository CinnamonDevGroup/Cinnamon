package minecraft

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/AngelFluffyOokami/Cinnamon/modules/core/commonutils"
	"github.com/AngelFluffyOokami/Cinnamon/modules/core/websocket"
	minecraftdb "github.com/AngelFluffyOokami/Cinnamon/modules/integrations/minecraft/database"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

var WebsocketHandler = map[string]func(data websocket.IncomingData, h *websocket.Hub){
	"minecraft": func(data websocket.IncomingData, h *websocket.Hub) {
		fmt.Print("msg received")
		var receivedData Data
		json.Unmarshal(data.RawData, &receivedData)
		client := <-h.Client

		if client.Authenticated {

			switch receivedData.DataType {
			case "playermessageevent":
				onPlayerMessage(receivedData)
			case "playerjoinevent":
				onPlayerJoin(receivedData.RawData)

			}

		} else {
			switch receivedData.DataType {
			case "authenticate":
				clientAuthenticate(client, h, receivedData.RawData)
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
				client.Send <- response

				close(client.Send)
				delete(h.Clients, client)

			}
		}
	},
}

func onPlayerMessage(data Data) {

	m := data.RawData
	config := <-commonutils.GetConfig

	s := <-commonutils.GetSession

	var chatMsg ChatMessage
	json.Unmarshal(m, &chatMsg)
	defer commonutils.RecoverPanic("")

	if chatMsg.Mention == "" {
		_, err := s.ChannelMessageSend("", chatMsg.Player+": "+chatMsg.Message)
		if err != nil {
			panic(err)
		}

	} else {
		message := discordgo.MessageReference{
			MessageID: chatMsg.Mention,
			ChannelID: "",
			GuildID:   "",
		}

		_, err := s.ChannelMessageSendReply("", chatMsg.Player+": "+chatMsg.Message, &message)
		if err != nil {
			panic(err)
		}
	}

	if config.Debugging {
		commonutils.LogEvent("Minecraft Player Message Event:\n"+chatMsg.Message+chatMsg.Player, commonutils.LogInfo)
	}
}

func clientAuthenticate(client *websocket.Client, h *websocket.Hub, responseData json.RawMessage) {
	s := <-commonutils.GetSession
	DB := <-commonutils.GetDB
	config := <-commonutils.GetConfig

	var authData Authenticate

	json.Unmarshal(responseData, &authData)
	defer commonutils.RecoverPanic(authData.DefaultChannel)
	connection := ConnectionStatus{
		AuthKey: authData.AuthKey,
	}

	server := minecraftdb.Minecraft{AuthKey: authData.AuthKey, GID: authData.GuildID}

	result := DB.First(&server)

	notexists := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if notexists {

		connection.Status = http.StatusUnauthorized
		connection.GID = ""

	} else {

		connection.Status = http.StatusOK
		connection.GID = server.GID
		client.AuthKey = authData.AuthKey
		client.Authenticated = true
		_, err := s.ChannelMessageSend(authData.DefaultChannel, "Minecraft server connected.")

		if err != nil {
			if config.Debugging {
				commonutils.LogEvent("ClientAuthenticate ChannelMessageSend Error Event: "+fmt.Sprint(err), commonutils.LogError)
			}
		} else {
			if config.Debugging {
				commonutils.LogEvent("ClientAuthenticate Client Register Event: "+client.Addr+" "+server.GID+" "+client.AuthKey, commonutils.LogInfo)
			}
		}
	}
	response, err := json.Marshal(connection)

	if err != nil {
		panic(err)
	}

	client.Send <- response
	if notexists {
		close(client.Send)
		delete(h.Clients, client)
	}

}

func onPlayerJoin(m []byte) {

	var playerJoined PlayerJoin
	json.Unmarshal(m, &playerJoined)
	defer commonutils.RecoverPanic("")

}
