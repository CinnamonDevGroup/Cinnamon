package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dop251/goja"
)

func InitCommands(vm *goja.Runtime) (commandsList string, err error) {
	path := "./modules"

	os.MkdirAll(path, os.ModePerm)
	files, err := os.ReadDir("./modules")
	if err != nil {
		return "", err
	}

	dirLen := len(files)

	modList := make([]string, 0)
	modContent := make([]string, 0)

	for x := 0; x < dirLen; x++ {
		if strings.HasSuffix(files[x].Name(), "mod.js") {
			modList = append(modList, files[x].Name())
			fileDir := "./modules/" + files[x].Name()
			modRead, err := os.ReadFile(fileDir)
			if err != nil {
				log.Fatal(err)
			}
			modContent = append(modContent, string(modRead))
		}
	}

	for x := 0; x < len(modContent); x++ {

		_, err := vm.RunString(modContent[x])
		slashDataReturn, ok := goja.AssertFunction(vm.Get("slashCommands"))

		if !ok {
			panic("Not a function")
		}

		slashData, err := slashDataReturn(goja.Undefined())

		slashSlice := slashData.Export()

		if err != nil {
			log.Fatal(err)
		}

		uwu, err := json.Marshal(slashSlice)
		fmt.Println(uwu)
	}

	fmt.Println(len(modList))
	return "", nil

}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("h")
}
