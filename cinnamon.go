package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	//	"gorm.io/gorm"
	"github.com/AngelFluffyOokami/Cinnamon/commands"
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-plugin/examples/grpc/shared"
	//	"gorm.io/driver/sqlite3"
)

func main() {

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
		fmt.Println(err)
		return
	}

	token := jsonFile["token"].(string)
	cookey = jsonFile["CookieKeyUUID"].(string)
	cookieVal = jsonFile["CookieValueUUID"].(string)
	if token == "inserttokenhere" {
		log.Fatal("No valid token has been detected, please insert bot token in the local configuration json file!")
	}

	addForeignPlugin(cookey, cookieVal)

	dgo, err := discordgo.New("Bot " + token)

	if err != nil {
		fmt.Println(err)
		return
	}

	schema := commands.MessageCreate
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

func addForeignPlugin(cookey string, cookieVal string) {

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.handshakeConfig,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command("~/./jumpstart"),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolGRPC},
	})
	defer client.Kill()

	rpcClient, err := client.Client()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	raw, err := rpcClient.Dispense("kv_grpc")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	kv := raw.(shared.KV)
	switch os.Args[0] {
	case "GET":
		result, err := kv.Get(os.Args[1])
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println(string(result))

	case "PUT":
		err := kv.Put(os.Args[1], []byte(os.Args[2]))
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

	default:
		fmt.Println("Only GET or PUT allowed", os.Args[0])
		os.Exit(1)
	}
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   cookey,
	MagicCookieValue: cookieVal,
}

var cookey string

var cookieVal string
