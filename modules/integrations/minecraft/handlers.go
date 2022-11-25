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
