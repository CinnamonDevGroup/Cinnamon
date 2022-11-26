package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	databaseHelper "github.com/AngelFluffyOokami/Cinnamon/modules/core/database"
	"github.com/AngelFluffyOokami/Cinnamon/modules/core/discord"
	coreserver "github.com/AngelFluffyOokami/Cinnamon/modules/core/server"
	"github.com/AngelFluffyOokami/Cinnamon/modules/integrations/minecraft"
	"github.com/bwmarrin/discordgo"
)

var s *discordgo.Session

func init() {

	//	if ./config.json exists, then:
	//	else if ./config.json does not exist, then:
	//	else if ./config.json is Schrödingers pet, then:
	if _, err := os.Stat("./config.json"); err == nil {
		// log to console: config file found
		log.Println("Configuration file found, loading configuration...")

	} else if errors.Is(err, os.ErrNotExist) {
		//	create a map with contents to include in the JSON file.
		config := map[string]interface{}{
			"token": "tokenhere",
		}
		//	Pretty JSON
		data, err := json.MarshalIndent(config, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		// 	create config.json, and write data to it
		_ = ioutil.WriteFile("config.json", data, 0644)
		log.Fatal("Configuration file not found! Creating new one. Please insert bot token into file before restarting")
		return
	} else {
		log.Fatal(err)
	}

	jsonFile, err := jsonParse()

	if err != nil {
		log.Fatal(err)
		return
	}

	disToken := jsonFile["token"].(string)

	DB := databaseHelper.Init()
	s = discord.Init(disToken)
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := minecraft.CommandsHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i, DB)
		}
	})
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := coreserver.CommandsHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i, DB)
		}
	})

	allCommands := append(minecraft.Commands, coreserver.Commands...)

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(allCommands))
	for i, v := range allCommands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", &v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	s.AddHandler(func(s *discordgo.Session, z *discordgo.GuildCreate) {
		coreserver.OnServerJoin(s, z, DB)
	})
	go minecraft.Init(s, DB)
}

func main() {

	fmt.Println("Bot is running")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("Removing commands...")

	commands, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		panic(err)
	}
	commandLen := len(commands)
	for x := 0; x < commandLen; x++ {
		err := s.ApplicationCommandDelete(s.State.User.ID, "", commands[x].ID)
		if err != nil {
			panic(err)
		}
	}
	s.Close()

}

func jsonParse() (map[string]interface{}, error) {
	jsonFile, err := os.Open("./config.json")
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	return result, err
}
