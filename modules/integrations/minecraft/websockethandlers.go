package minecraft

import (
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func onPlayerMessage(m []byte, s *discordgo.Session) {

	var chatMsg chatMessage
	json.Unmarshal(m, &chatMsg)

	fmt.Print(chatMsg.Mention + chatMsg.Message + chatMsg.Player)
}

func onClientConnect(multidb *Multidbmodel, client *Client) {

	clientAuthID := client.AuthID
	fmt.Println(clientAuthID)

}

func onPlayerJoin(m []byte, s *discordgo.Session) {

	var playerJoined playerJoin
	json.Unmarshal(m, &playerJoined)

}
