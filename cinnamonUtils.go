package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/AngelFluffyOokami/Cinnamon/modules/core/commonutils"
	coreserver "github.com/AngelFluffyOokami/Cinnamon/modules/core/server"
	"github.com/bwmarrin/discordgo"
)

const (
	token           = "tokenhere"
	adminserver     = "adminserveridhere"
	adminchannel    = "adminchannelidhere"
	infochannel     = "infochannelidhere"
	warnchannel     = "warnchannelidhere"
	errchannel      = "errorchannelidhere"
	updatechannel   = "updatechannelidhere"
	feedbackchannel = "feedbackchannelidhere"
)

func ReadJSON(filename string) (commonutils.Data, error) {
	config := <-commonutils.GetConfig
	data, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func initDiscordHandlers() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := allCommandHandlers[i.ApplicationCommandData().Name]; ok {
			fmt.Println(i.ApplicationCommandData().Name)
			h(i)
		}
	})

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
		coreserver.OnServerJoin(z)
	})
}

func exit() {
	fmt.Println("Removing commands...")

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

func beautifyJSONFile(filename string) {
	// Open the given file
	jsonFile, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	// Unmarshal the contents of the file into an interface
	var jsonData interface{}
	json.Unmarshal(jsonFile, &jsonData)

	// Marshal the data into a byte array with indentation
	prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	// Write the byte array back to the file
	err = os.WriteFile(filename, prettyJSON, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

// CreateOrUpdateJSON creates or updates a JSON file with two keys
func CreateOrUpdateJSON(file string) error {
	// Read the existing file
	data := commonutils.Data{}
	bytes, err := os.ReadFile(file)
	if err != nil {
		// File does not exist, create it
		data = commonutils.Data{
			Token:           token,
			AdminServer:     adminserver,
			AdminChannel:    adminchannel,
			InfoChannel:     infochannel,
			WarnChannel:     warnchannel,
			ErrChannel:      errchannel,
			UpdateChannel:   updatechannel,
			FeedbackChannel: feedbackchannel,
		}
	} else {
		// File exists, parse it
		if err := json.Unmarshal(bytes, &data); err != nil {
			return fmt.Errorf("failed to parse existing file: %v", err)
		}
		// Check if key1 or key2 are missing
		if data.Token == "" {
			data.Token = token
		}
		if data.AdminServer == "" {
			data.AdminServer = adminserver
		}
		if data.AdminChannel == "" {
			data.AdminChannel = adminchannel
		}
		if data.InfoChannel == "" {
			data.InfoChannel = infochannel
		}
		if data.WarnChannel == "" {
			data.WarnChannel = warnchannel
		}
		if data.ErrChannel == "" {
			data.ErrChannel = errchannel
		}
		if data.UpdateChannel == "" {
			data.UpdateChannel = updatechannel
		}
		if data.FeedbackChannel == "" {
			data.FeedbackChannel = feedbackchannel
		}
	}

	// Marshal the data back to JSON
	bytes, err = json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	// Write the data back to the file
	if err := os.WriteFile(file, bytes, 0644); err != nil {
		return fmt.Errorf("failed to write data to file: %v", err)
	}

	return nil
}
