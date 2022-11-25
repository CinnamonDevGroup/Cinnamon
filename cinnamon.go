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

	cinnamonModel "github.com/AngelFluffyOokami/Cinnamon/modules/database/cinnamon"
	guildModel "github.com/AngelFluffyOokami/Cinnamon/modules/database/guild"
	userModel "github.com/AngelFluffyOokami/Cinnamon/modules/database/user"
	"github.com/AngelFluffyOokami/Cinnamon/modules/integrations/minecraft"
	"github.com/bwmarrin/discordgo"
	"github.com/glebarez/sqlite"
	"github.com/zmb3/spotify"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
	"gorm.io/gorm"
)

var Client spotify.Client
var s *discordgo.Session
var guildDB *gorm.DB
var cinnamonDB *gorm.DB
var userDB *gorm.DB
var err error

type multidbmodel struct {
	user     *gorm.DB
	guild    *gorm.DB
	cinnamon *gorm.DB
}

var multidb multidbmodel

func init() {

	guildDB, err = gorm.Open(sqlite.Open("database/guilds.db"), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}
	userDB, err = gorm.Open(sqlite.Open("database/users.db"), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}
	cinnamonDB, err = gorm.Open(sqlite.Open("database/cinnamon.db"), &gorm.Config{})

	guildDB.AutoMigrate(&guildModel.Guild{})

	userDB.AutoMigrate(&userModel.User{})

	cinnamonDB.AutoMigrate(&cinnamonModel.Cinnamon{})

	multidb.user = userDB

	multidb.guild = guildDB

	multidb.cinnamon = cinnamonDB

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

	err = s.Open()

	if err != nil {
		log.Fatal(err)
	}
	go minecraft.Run(s)
}

func main() {

	fmt.Println("Bot is running")

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
