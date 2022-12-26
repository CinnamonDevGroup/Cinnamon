package commonutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime"
	"time"

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
	key := babbler.Babble()
	return key
}

func initializeServer(GID string) {
	defer RecoverPanic("")

	s := <-GetSession
	DB := <-GetDB

	var JoinedAt []int64

	JoinedAt = append(JoinedAt, time.Now().Unix())

	var MemberCount int
	guildCheck, err := s.Guild(GID)
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

	servername := GetGuildName(GID)

	ownername := GetGuildOwnerName(GID)

	message := "Server Join Event: " + servername + " " + ownername + "\n"
	LogEvent(message, LogInfo)

	result := DB.Create(&guild)
	fmt.Print(result.Error)

}

const (
	LogError    = "ERR"
	LogWarning  = "WARN"
	LogInfo     = "INFO"
	LogUpdate   = "UPDATE"
	LogFeedback = "FEEDBACK"
)

func GetGuildName(GID string) string {
	s := <-GetSession
	guild, err := s.Guild(GID)
	if err != nil {
		return "Undefined. " + GID
	} else {
		return guild.Name + " " + GID
	}
}

func GetGuildOwnerName(GID string) string {
	s := <-GetSession
	g, err := s.Guild(GID)
	if err != nil {
		return "Undefined."
	}

	user, err := s.User(g.OwnerID)
	if err != nil {
		return "Undefined. " + g.OwnerID
	} else {
		return user.Username + "#" + user.Discriminator + " " + g.OwnerID
	}
}

func LogEvent(message string, level string) {

	config := <-GetConfig
	s := <-GetSession
	// create the log entry
	entry := LogEntry{
		Time:    time.Now(),
		Message: message,
		Level:   level,
	}

	// open the log file
	logFile, err := os.OpenFile("opossum.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	// marshal the log entry
	entryJSON, err := json.MarshalIndent(entry, "\n", "\n")
	if err != nil {
		panic(err)
	}

	// write the log entry to the file
	if _, err := logFile.Write(entryJSON); err != nil {
		panic(err)
	}
	if err := logFile.Sync(); err != nil {
		panic(err)
	}
	switch level {
	case "INFO":
		s.ChannelMessageSend(config.InfoChannel, entry.Level+": \n"+entry.Message+"\n"+fmt.Sprint(entry.Time))
	case "WARN":
		s.ChannelMessageSend(config.WarnChannel, entry.Level+": \n"+entry.Message+"\n"+fmt.Sprint(entry.Time))
	case "ERR":
		s.ChannelMessageSend(config.ErrChannel, entry.Level+": \n"+entry.Message+"\n"+fmt.Sprint(entry.Time))
	case "UPDATE":
		s.ChannelMessageSend(config.UpdateChannel, entry.Level+": \n"+entry.Message+"\n"+fmt.Sprint(entry.Time))
	case "FEEDBACK":
		s.ChannelMessageSend(config.FeedbackChannel, entry.Level+": \n"+entry.Message)
	}
}

var SetConfig = make(chan Data)
var GetConfig = make(chan Data)

func Config() {
	var config Data
	for {
		select {
		case config = <-SetConfig:
		case GetConfig <- config:
		}
	}

}

var GetSession = make(chan *discordgo.Session)
var SetSession = make(chan *discordgo.Session)

func Session() {
	var session *discordgo.Session
	for {
		select {
		case session = <-SetSession:
		case GetSession <- session:
		}
	}

}

var GetDB = make(chan *gorm.DB)
var SetDB = make(chan *gorm.DB)

func DB() {
	var DB *gorm.DB
	for {
		select {
		case DB = <-SetDB:
		case GetDB <- DB:
		}
	}
}

func RecoverPanic(channelID string) {

	if r := recover(); r != nil {

		s := <-GetSession

		// get the stack trace of the panic
		tempbuf := make([]byte, 10000)
		buflength := runtime.Stack(tempbuf, false)
		var buf []byte
		if buflength >= 1900 {
			buf = make([]byte, 1900)
		} else {
			buf = make([]byte, buflength)
		}
		runtime.Stack(buf, false)

		LogEvent(fmt.Sprintf("Recovering from panic: %v\n Stack trace: %s", r, buf), "ERR")
		if channelID != "" {
			s.ChannelMessageSend(channelID, "Error processing command.\nBug report sent to developers.")
		}

	}

}

func CheckGuildExists(GID string) {
	DB := <-GetDB
	guild := coredb.Guild{GID: GID}

	result := DB.First(&guild)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		initializeServer(GID)

	} else {
		return
	}
}
