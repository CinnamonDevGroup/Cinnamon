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

	//	"gorm.io/gorm"
	"github.com/AngelFluffyOokami/Cinnamon/commands"
	"github.com/bwmarrin/discordgo"
	//	"gorm.io/driver/sqlite3"
)

type Rankings struct {
	Token string `json:"token"`
}

func main() {

	schema := commands.MessageCreate

	if _, err := os.Stat("./config.json"); err == nil {
		log.Println("Configuration file found, loading configuration...")

	} else if errors.Is(err, os.ErrNotExist) {
		config := map[string]interface{}{
			"token": "inserttokenhere",
		}
		data, err := json.MarshalIndent(config, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		_ = ioutil.WriteFile("config.json", data, 0644)
		log.Fatal("Configuration file not found! Creating new one. Please insert bot token into file before restarting")
		return
	} else {
		log.Fatal(err)
	}

	jsonFile, err := jsonParse()

	if err != nil {
		fmt.Println(err)
		return
	}

	token := jsonFile["token"].(string)
	if token == "inserttokenhere" {
		log.Fatal("No valid token has been detected, please insert bot token in the local configuration json file!")
	}

	dgo, err := discordgo.New("Bot " + token)

	if err != nil {
		fmt.Println(err)
		return
	}

	dgo.AddHandler(schema)
	dgo.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	err = dgo.Open()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Bot is running")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dgo.Close()

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
