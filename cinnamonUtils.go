package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/CinnamonDevGroup/Cinnamon/modules/core/common"
	"github.com/CinnamonDevGroup/Cinnamon/modules/core/core_handlers"
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

// Function reads json file, returning variable of type config.Data
func ReadJSON(filename string) (common.Data, error) {
	//	Reads file and saves []byte to variable data, then checks if there was an error, if error, then return nil config and error.
	data, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	//	Unmarshals []byte of type json into variable config of type config.Data, if error, then return nil config, and error.
	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	// 	If no error, return config, and nil error.
	return config, nil
}

// Registers commands, and adds command handling functions to discord session after discord session is opened.
func initDiscordHandlers() {
	//	Adds handler of type i *discordgo.InteractionCreate that selects the appropriate handler over a map of command handlers depending on the application
	s.AddHandler(func(i *discordgo.InteractionCreate) {
		if h, ok := allCommandHandlers[i.ApplicationCommandData().Name]; ok {
			fmt.Println(i.ApplicationCommandData().Name)
			h(i)
		}
	})

	log.Println("Adding commands...")
	// registers commands with the discord session.
	registeredCommands := make([]*discordgo.ApplicationCommand, len(allCommands))
	for i, v := range allCommands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", &v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	// adds core handlers
	s.AddHandler(func(s *discordgo.Session, z *discordgo.GuildCreate) {
		core_handlers.OnServerJoin(z)
	})
}

// Clean up session by removing commands and closing the session afterwards
func exit() {
	fmt.Println("Removing commands...")

	// Get a list of all the commands registered within.
	commands, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		panic(err)
	}

	// iterates through all commands and deletes them.
	for x := 0; x < len(commands); x++ {
		err := s.ApplicationCommandDelete(s.State.User.ID, "", commands[x].ID)
		if err != nil {
			panic(err)
		}
	}
	// closes discord session.
	s.Close()
}

// beautifies a json file
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
	data := common.Data{}
	bytes, err := os.ReadFile(file)
	if err != nil {
		// File does not exist, create it
		data = common.Data{
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
		// Check if keys are missing.
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
