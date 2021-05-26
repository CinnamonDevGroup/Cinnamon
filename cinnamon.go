package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/AngelFluffyOokami/Cinnamon/commandshandler/commandsschemas"

	//	"gorm.io/gorm"
	"github.com/bwmarrin/discordgo"
	//	"gorm.io/driver/sqlite3"
)


func main() {


	jsonFile, err := jsonParse()
	
	if err != nil {
		fmt.Println(err)
		return
	}

	token := jsonFile["Token"].(string)
	dgo, err := discordgo.New("Bot " + token)

	if err != nil {
		fmt.Println(err)
		return
	}

	dgo.AddHandler(commandsschemas.MessageCreate)
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
