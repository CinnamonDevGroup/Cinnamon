package minecraft

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	minecraftdb "github.com/AngelFluffyOokami/Cinnamon/modules/integrations/minecraft/database"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

func onPlayerMessage(m []byte, s *discordgo.Session) {

	var chatMsg chatMessage
	json.Unmarshal(m, &chatMsg)

	if chatMsg.Mention == "" {
		_, err := s.ChannelMessageSend(chatMsg.Channel, chatMsg.Player+": "+chatMsg.Message)
		if err != nil {
			panic(err)
		}

	} else {
		message := discordgo.MessageReference{
			MessageID: chatMsg.Mention,
			ChannelID: chatMsg.Channel,
			GuildID:   chatMsg.Guild,
		}

		_, err := s.ChannelMessageSendReply(chatMsg.Channel, chatMsg.Player+": "+chatMsg.Message, &message)
		if err != nil {
			panic(err)
		}
	}

	fmt.Print(chatMsg.Mention + chatMsg.Message + chatMsg.Player)
}

func clientAuthenticate(client *Client, DB *gorm.DB, s *discordgo.Session, h *Hub, responseData json.RawMessage) {

	var authData authenticate

	json.Unmarshal(responseData, &authData)
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
		fmt.Println(err)
	}
	response, err := json.Marshal(connection)

	if err != nil {
		panic(err)
	}

	clientAuthKey := client.AuthKey
	fmt.Println(clientAuthKey)
	client.send <- response
	if notexists {
		close(client.send)
		delete(h.clients, client)
	}

}

func onPlayerJoin(m []byte, s *discordgo.Session) {

	var playerJoined playerJoin
	json.Unmarshal(m, &playerJoined)

}
