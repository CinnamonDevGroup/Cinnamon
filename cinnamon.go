package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/CinnamonDevGroup/Cinnamon/modules/core/common"
	"github.com/CinnamonDevGroup/Cinnamon/modules/core/core_handlers"
	"github.com/CinnamonDevGroup/Cinnamon/modules/core/database"
	discord_session "github.com/CinnamonDevGroup/Cinnamon/modules/core/discord"
	"github.com/CinnamonDevGroup/Cinnamon/modules/core/websocket"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

var s *discordgo.Session
var DB *gorm.DB

/*
* 	These maps are created for any module added using a build tag to add their functions upon init()
* 	These maps get populated by the init() functions of each module added via a build tag.
*
* 	The slice containing discordgo.ApplicationCommands gets appended to by the same modules via build tags to add discord application commands.
*
* 	Function slice DBMigrate gets appended to by every module that needs to initialize its own database table.
 */
var allCommandHandlers = make(map[string]func(i *discordgo.InteractionCreate))
var allWebsocketHandlers = make(map[string]func(receivedData websocket.IncomingData, h *websocket.Hub))
var allCommands []discordgo.ApplicationCommand
var DBMigrate []func()

/*
*	Variable config created to populate it with the info fetched from the configuration json file.
*
*	Bool channel off gets created to be able to signal the bot to shutdown at will.
*
*	Variable error created to be able to start the discord session without overwriting the global discord session variable.
 */
var config common.Data
var off = make(chan bool)
var err error

func init() {

	/*
	*	CreateOrUpdateJSON() creates a json configuration file if not exists, if exists and doesn't have all the configuration options,
	*	it updates the file to contain the missing config options, leaving the rest untouched in their state.
	*
	*	beautifyJSONFile() beautifies the configuration file created or updated by the CreateOrUpdateJSON file.
	*
	*	ReadJSON() reads the configuration file and saves it to the global variable config.
	 */
	CreateOrUpdateJSON("config.json")
	beautifyJSONFile("config.json")
	config, err = ReadJSON("config.json")
	if err != nil {
		panic(err)
	}

	/*
	*	Variable config saves its value to common.Config for other modules to access without causing circular dependency import
	* 	by attempting to access config from main package.
	 */
	common.Config = config

	/*
	*	Database gets initialized, returning the DB engine to the variable DB, which then gets written to common.DB, for other
	*	modules to access without causing circular dependency import by attempting to use DB straight from main package.
	*
	*	Discord session gets initialized then returns the session to variable s, which proceeds to get written to common.Session for other
	* 	modules to access without causing circular dependency import by attempting to use session from main package.
	 */
	DB = database.InitDB()
	common.DB = DB
	s = discord_session.InitSession(config.Token)
	common.Session = s

	allCommands = append(allCommands, core_handlers.Commands...)
	for k, v := range core_handlers.CommandsHandlers {
		allCommandHandlers[k] = v
	}
}

func main() {
	websocket.WebsocketHandlers = allWebsocketHandlers
	go websocket.InitWebsocket()
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
