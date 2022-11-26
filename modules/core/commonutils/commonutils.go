package commonutils

import (
	"errors"
	"time"

	databaseHelper "github.com/AngelFluffyOokami/Cinnamon/modules/core/database"
	coredb "github.com/AngelFluffyOokami/Cinnamon/modules/core/database/core"
	"github.com/bwmarrin/discordgo"
	"github.com/tjarratt/babble"
	"gorm.io/gorm"
)

func BabbleWords() string {
	var wordlist []string
	wordlist = append(wordlist, "grudge", "linear", "burial", "latest", "screen", "desert", "expose", "endure", "estate", "master", "refund", "throat", "effort", "pepper", "budget", "revive", "breast", "school", "flower", "ladder", "chorus", "wonder", "cheese", "sticky", "spread", "tumble", "vacuum", "flavor", "suntan", "mutter", "center", "punish", "resort", "hunter", "galaxy", "charge", "depend", "cotton", "shiver", "afford", "agenda", "timber", "morale", "behave", "camera", "expand", "carbon", "dollar", "latest", "mature", "mobile", "injury", "ensure", "barrel", "finish", "rhythm", "crutch", "museum", "lesson", "follow", "please", "safety", "modest", "remind", "reader", "demand", "ethics", "pledge", "accept", "ballot", "doctor", "gutter", "planet", "launch", "makeup", "freeze", "acquit", "colony", "rescue", "defend", "facade", "vision", "honest", "retire", "arrest", "banner", "thesis", "weight", "turkey", "worker", "column", "ignite", "facade", "ribbon", "bloody", "sacred", "inside", "dilute", "gallon", "theory", "behead", "proper", "chance", "single", "object", "temple", "modest", "likely", "adjust", "pastel", "attack", "market", "bishop", "belong", "effort", "rotate", "senior", "infect", "locate", "secure", "earwax", "normal", "flower", "prayer", "endure", "injury", "avenue", "family", "desert", "packet", "series", "tiptoe", "tumble", "harass", "spider", "output", "mutter", "church", "glance", "throne", "salmon", "option", "apathy", "cancer", "labour", "stroke", "dinner", "lounge", "gallon", "mobile", "bubble", "trance", "matrix", "ground", "escape", "defeat", "effect", "acquit", "square", "bitter", "excuse", "review", "normal", "formal", "player", "quaint", "belief", "critic", "accent", "empire", "junior", "lesson", "tongue", "voyage", "basket", "launch", "mosaic", "column", "margin", "source", "spirit", "cherry", "height", "bother", "deadly", "marble", "virtue", "devote", "mosque", "morale", "likely", "branch", "offend", "family", "script", "medium", "course", "theory", "weight", "winner")
	babbler := babble.NewBabbler()
	babbler.Count = 6
	babbler.Words = wordlist
	return babbler.Babble()
}

func initializeServer(GID string, DB databaseHelper.DBstruct, s *discordgo.Session) {

	var JoinedAt []int64

	JoinedAt = append(JoinedAt, time.Now().Unix())

	var MemberCount int
	guildCheck, err := s.State.Guild(GID)
	if err != nil {
		MemberCount = 0
	} else {
		MemberCount = guildCheck.MemberCount
	}

	var messages []coredb.Message

	messages = append(messages, coredb.Message{
		MessageCount: 0,
		TimeCount:    time.Now().Unix(),
	})

	guild := coredb.Guild{
		GID:     GID,
		AuthKey: BabbleWords(),
		Joined:  time.Now().Unix(),
		About: coredb.Information{
			JoinedAt:     JoinedAt,
			UserAmount:   MemberCount,
			MessageCount: messages,
		},
	}

	DB.Guilds.Create(&guild)

}

func CheckServerExists(GID string, DB databaseHelper.DBstruct, s *discordgo.Session) {
	guild := coredb.Guild{GID: GID}

	result := DB.Guilds.First(&guild)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		initializeServer(GID, DB, s)

	} else {
		return
	}
}
