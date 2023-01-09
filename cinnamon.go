package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/CinnamonDevGroup/Cinnamon/modules/core/common"
	"github.com/CinnamonDevGroup/Cinnamon/modules/core/core_handlers"
	"github.com/CinnamonDevGroup/Cinnamon/modules/core/database"
	discord_client "github.com/CinnamonDevGroup/Cinnamon/modules/core/discord"
	"github.com/CinnamonDevGroup/Cinnamon/modules/core/websocket"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

var s *discordgo.Session
var DB *gorm.DB

var allCommandHandlers = make(map[string]func(i *discordgo.InteractionCreate))
var allWebsocketHandlers = make(map[string]func(receivedData websocket.IncomingData, h *websocket.Hub))
var allCommands []discordgo.ApplicationCommand
var DBMigrate []func()
var config common.Data
var off = make(chan bool)
var err error

func init() {

	CreateOrUpdateJSON("config.json")
	beautifyJSONFile("config.json")
	config, err = ReadJSON("config.json")
	common.Config = config
	if err != nil {
		panic(err)
	}
	DB = database.Init()
	common.DB = DB
	s = discord_client.Init(config.Token)

	common.Session = s

	allCommands = append(allCommands, core_handlers.Commands...)
	for k, v := range core_handlers.CommandsHandlers {
		allCommandHandlers[k] = v
	}
}

func main() {
	websocket.WebsocketHandlers = allWebsocketHandlers
	go websocket.Init()
	for _, x := range DBMigrate {
		x()
	}
	initDiscordHandlers()

	fmt.Println("Bot is running")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	select {
	case <-sc:
	case <-off:
	}
	exit()

}
