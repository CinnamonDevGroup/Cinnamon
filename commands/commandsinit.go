package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/dop251/goja"
)

func InitCommands(vm *goja.Runtime) (commandsList string, err error) {
	fmt.Print("1")
	return

}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("h")
}
