package minecraft_websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/AngelFluffyOokami/Cinnamon/modules/core/commonutils"
	coredb "github.com/AngelFluffyOokami/Cinnamon/modules/core/database/core"
	"github.com/AngelFluffyOokami/Cinnamon/modules/core/websocket"
	minecraft_api_v1 "github.com/AngelFluffyOokami/Cinnamon/modules/integrations/minecraft/API/v1"
	minecraft_db "github.com/AngelFluffyOokami/Cinnamon/modules/integrations/minecraft/database"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

var WebsocketHandler = map[string]func(data websocket.IncomingData, h *websocket.Hub){
	"minecraft": func(data websocket.IncomingData, h *websocket.Hub) {

		fmt.Print("msg received")
		var receivedData minecraft_api_v1.Data
		json.Unmarshal(data.RawData, &receivedData)

		if data.Client.Authenticated {

			var cache websocket.ClientCache

			for _, x := range websocket.ClientsCache {
				if x.UUID == data.Client.UUID {
					cache = x
				}
			}
			switch receivedData.DataType {
			case "playermessageevent":
				onPlayerMessage(receivedData)
			case "playerjoinevent":
				onPlayerJoin(receivedData.RawData, cache)
			}

		} else {
			switch receivedData.DataType {
			case "authenticate":
				clientAuthenticate(data.Client, h, receivedData.RawData)
			default:
				connection := minecraft_api_v1.ConnectionStatus{
					AuthKey: "",
					GID:     "",
					Status:  http.StatusUnauthorized,
				}
				response, err := json.Marshal(connection)
				if err != nil {
					panic(err)
				}
				data.Client.Send <- response

				close(data.Client.Send)
				delete(h.Clients, data.Client)

			}
		}
	},
}

func onPlayerMessage(data minecraft_api_v1.Data) {

	m := data.RawData
	config := commonutils.Config

	s := commonutils.Session

	var chatMsg minecraft_api_v1.ChatMessage
	json.Unmarshal(m, &chatMsg)
	defer commonutils.RecoverPanic("")

	if chatMsg.Mention == "" {
		_, err := s.ChannelMessageSend("", ": "+chatMsg.Message)
		if err != nil {
			panic(err)
		}

	} else {
		message := discordgo.MessageReference{
			MessageID: chatMsg.Mention,
			ChannelID: "",
			GuildID:   "",
		}

		_, err := s.ChannelMessageSendReply("", ": "+chatMsg.Message, &message)
		if err != nil {
			panic(err)
		}
	}

	if config.Debugging {
		commonutils.LogEvent("Minecraft Player Message Event:\n"+chatMsg.Message+chatMsg.UUID, commonutils.LogInfo)
	}
}

func clientAuthenticate(client *websocket.Client, h *websocket.Hub, responseData json.RawMessage) {
	s := commonutils.Session
	DB := commonutils.DB
	config := commonutils.Config

	var authData minecraft_api_v1.Authenticate

	json.Unmarshal(responseData, &authData)
	defer commonutils.RecoverPanic(authData.DefaultChannel)
	connection := minecraft_api_v1.ConnectionStatus{
		AuthKey: authData.AuthKey,
	}

	server := minecraft_db.Minecraft{AuthKey: authData.AuthKey, GID: authData.GuildID}

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
		newClient := websocket.ClientCache{
			Client:         client,
			UUID:           client.UUID,
			GID:            authData.GuildID,
			DefaultChannel: authData.DefaultChannel,
			AuthKey:        authData.AuthKey,
			Service:        "minecraft",
		}
		websocket.ClientsCache = append(websocket.ClientsCache, newClient)
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

func checkPlayerAuth(UUID string, cache websocket.ClientCache) (bool, coredb.Service) {
	DB := commonutils.DB

	currentUserService := coredb.Service{
		Service: "minecraft",
		UUID:    UUID,
	}
	result := DB.First(&currentUserService)
	return !(errors.Is(result.Error, gorm.ErrRecordNotFound)), currentUserService

}

func KickNewPlayerAuth(UUID string, cache websocket.ClientCache) {

	DB := commonutils.DB

	newUser := coredb.Service{
		Service: "minecraft",
		UUID:    UUID,
	}

	newAuthKey := commonutils.BabbleWords()

	newUser.AuthKey = newAuthKey

	DB.Save(&newUser)

	kickAuth := minecraft_api_v1.KickForAuth{
		UUID:    UUID,
		AuthKey: newUser.AuthKey,
	}
	AuthKick, err := json.Marshal(kickAuth)
	if err != nil {
		if commonutils.Config.Debugging {
			commonutils.LogEvent("JSON Marshal Error Event: "+fmt.Sprint(err), commonutils.LogError)
		}
		return
	}
	OutData := minecraft_api_v1.OutboundData{
		DataType: minecraft_api_v1.AuthKickEvent,
		RawData:  AuthKick,
		API:      cache.Client.APIVersion,
	}
	response, err := json.Marshal(OutData)
	if err != nil {
		if commonutils.Config.Debugging {
			commonutils.LogEvent("JSON Marshal Error Event: "+fmt.Sprint(err), commonutils.LogError)
		}
		return
	}
	cache.Client.Send <- response

}

func kickPlayerAuth(UUID string, Cache websocket.ClientCache) {
	DB := commonutils.DB

	currentUser := coredb.Service{
		Service: "minecraft",
		UUID:    UUID,
	}
	DB.First(&currentUser)

	KickAuth := minecraft_api_v1.KickForAuth{
		UUID:    UUID,
		AuthKey: currentUser.AuthKey,
	}
	AuthKick, err := json.Marshal(KickAuth)
	if err != nil {
		if commonutils.Config.Debugging {
			commonutils.LogEvent("JSON Marshal Error Event: "+fmt.Sprint(err), commonutils.LogError)
		}
		return
	}

	OutData := minecraft_api_v1.OutboundData{
		DataType: minecraft_api_v1.AuthKickEvent,
		RawData:  AuthKick,
		API:      Cache.Client.APIVersion,
	}
	response, err := json.Marshal(OutData)
	if err != nil {
		if commonutils.Config.Debugging {
			commonutils.LogEvent("JSON Marshal Error Event: "+fmt.Sprint(err), commonutils.LogError)
		}
		return
	}

	Cache.Client.Send <- response
}

func DecideAuth(UUID string, Cache websocket.ClientCache) {
	DB := commonutils.DB
	s := commonutils.Session

	currentUserService := coredb.Service{
		Service: "minecraft",
		UUID:    UUID,
	}
	DB.First(&currentUserService)

	guilds, err := s.UserGuilds(100, "", "")

	if err != nil {
		commonutils.LogEvent("UserGuilds Error Event: "+fmt.Sprint(err), commonutils.LogError)
		return
	}

	for {
		if len(guilds) == 100 {
			newGuilds, err := s.UserGuilds(100, "", guilds[99].ID)
			if err != nil {
				commonutils.LogEvent("UserGuilds Error Event: "+fmt.Sprint(err), commonutils.LogError)
				break
			}
			for {
				guilds = append(guilds, newGuilds...)
				if len(newGuilds) == 100 {
					oldGuild := newGuilds[99]
					newGuilds = nil
					newGuilds, err = s.UserGuilds(100, "", oldGuild.ID)
					if err != nil {
						commonutils.LogEvent("UserGuilds Error Event: "+fmt.Sprint(err), commonutils.LogError)
						break
					}
				} else {
					break
				}
			}
		} else {
			break
		}
	}

	found := false
	for _, x := range guilds {
		if x.ID == Cache.GID {
			found = true
		}
	}

	user, _ := s.User(currentUserService.UID)

	if found {
		authSuccess(UUID, user.Username, Cache)
	} else {
		authFail(UUID, Cache)
	}

}

func authFail(UUID string, Cache websocket.ClientCache) {
	FailAuth := minecraft_api_v1.KickForNotOnServer{
		UUID: UUID,
	}
	data, err := json.Marshal(FailAuth)
	if err != nil {
		if commonutils.Config.Debugging {
			commonutils.LogEvent("JSON Marshal Error Event: "+fmt.Sprint(err), commonutils.LogError)
		}
		return
	}

	responsedata := minecraft_api_v1.OutboundData{
		DataType: minecraft_api_v1.NotFoundEvent,
		RawData:  data,
		API:      Cache.Client.APIVersion,
	}

	response, err := json.Marshal(responsedata)
	if err != nil {
		if commonutils.Config.Debugging {
			commonutils.LogEvent("JSON Marshal Error Event: "+fmt.Sprint(err), commonutils.LogError)
		}
		return
	}
	Cache.Client.Send <- response

}

func authSuccess(UUID string, username string, Cache websocket.ClientCache) {
	AuthSuccess := minecraft_api_v1.PlayerAuthSuccessful{
		UUID:     UUID,
		Username: username,
	}

	data, err := json.Marshal(AuthSuccess)
	if err != nil {
		if commonutils.Config.Debugging {
			commonutils.LogEvent("JSON Marshal Error Event: "+fmt.Sprint(err), commonutils.LogError)
		}
		return
	}

	responsedata := minecraft_api_v1.OutboundData{
		DataType: minecraft_api_v1.PlayerAuthEvent,
		RawData:  data,
		API:      Cache.Client.APIVersion,
	}
	response, err := json.Marshal(responsedata)
	if err != nil {
		if commonutils.Config.Debugging {
			commonutils.LogEvent("JSON Marshal Error Event: "+fmt.Sprint(err), commonutils.LogError)
		}
		return
	}

	Cache.Client.Send <- response
}

func onPlayerJoin(m []byte, Cache websocket.ClientCache) {

	var playerJoined minecraft_api_v1.PlayerJoin
	json.Unmarshal(m, &playerJoined)
	defer commonutils.RecoverPanic(Cache.DefaultChannel)

	exists, user := checkPlayerAuth(playerJoined.UUID, Cache)
	if !exists {
		KickNewPlayerAuth(playerJoined.UUID, Cache)
		return
	}
	if user.UID == "" {
		kickPlayerAuth(playerJoined.UUID, Cache)
		return
	}
	DecideAuth(playerJoined.UUID, Cache)
}
