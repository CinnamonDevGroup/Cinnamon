package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/AngelFluffyOokami/Cinnamon/modules/core/commonutils"
	databaseHelper "github.com/AngelFluffyOokami/Cinnamon/modules/core/database"
	"github.com/AngelFluffyOokami/Cinnamon/modules/core/discord"
	coreserver "github.com/AngelFluffyOokami/Cinnamon/modules/core/server"
	"github.com/AngelFluffyOokami/Cinnamon/modules/core/websocket"
	"github.com/AngelFluffyOokami/Cinnamon/modules/integrations/personalization/personalization"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

var s *discordgo.Session
var DB *gorm.DB

var allCommandHandlers = make(map[string]func(i *discordgo.InteractionCreate))
var allWebsocketHandlers = make(map[string]func(receivedData websocket.IncomingData, h *websocket.Hub))
var allCommands []discordgo.ApplicationCommand
var DBMigrate []func()
var config commonutils.Data
var off = make(chan bool)
var err error

func init() {

	go personalization.InitPersonalization()
	CreateOrUpdateJSON("config.json")
	beautifyJSONFile("config.json")
	config, err = ReadJSON("config.json")
	commonutils.Config = config
	if err != nil {
		panic(err)
	}
	DB = databaseHelper.Init()
	commonutils.DB = DB
	s = discord.Init(config.Token)

	commonutils.Session = s

	allCommands = append(allCommands, coreserver.Commands...)
	for k, v := range coreserver.CommandsHandlers {
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
