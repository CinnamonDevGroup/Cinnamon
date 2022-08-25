package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	//	"gorm.io/gorm"
	"github.com/AngelFluffyOokami/Cinnamon/voice"
	"github.com/bwmarrin/discordgo"
	"github.com/zmb3/spotify"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
	//	"gorm.io/driver/sqlite3"
)

var Client spotify.Client
var s *discordgo.Session
var control chan bool

var servers voice.ServersStruct

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
			"token":        "tokenhere",
			"clientID":     "clientIDhere",
			"clientSecret": "clientSecrethere",
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
	spotClientID := jsonFile["clientID"].(string)
	spotClientSecret := jsonFile["clientSecret"].(string)
	ctx := context.Background()
	config := &clientcredentials.Config{
		ClientID:     spotClientID,
		ClientSecret: spotClientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}

	token, err := config.Token(ctx)

	httpClient := spotifyauth.New().Client(ctx, token)

	Client = spotify.NewClient(httpClient)

	if err != nil {
		log.Fatal(err)
	}
	s, err = discordgo.New("Bot " + disToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Interaction.Type == 2 {
			if h, ok := voice.CommandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i, Client, &servers)
			}
		} else if i.Interaction.Type == 3 {
			voice.OnInteractionResponse(s, i, &servers, control)
		}

	})

	err = s.Open()

	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	fmt.Println("Bot is running")

	registeredCommands := make([]*discordgo.ApplicationCommand, len(voice.Commands))

	for i, v := range voice.Commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", &v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

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
