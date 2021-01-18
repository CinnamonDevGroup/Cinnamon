package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/AngelFluffyOokami/dbase"
	"github.com/bwmarrin/discordgo"
)

var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot token")
	flag.Parse()
}

func main() {
	dgo, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error encountered while establishing Discord session, ", err)
		return
	} 
	
	dgo.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	err = dgo.Open()
	if err != nil {
		fmt.Println("Error encountered while opening connection, ", err)
		return
	}

	fmt.Println("Bot is now running.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc 

	dgo.Close()
}