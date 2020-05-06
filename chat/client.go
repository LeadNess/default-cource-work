package main

import (
	"./tui"
	"log"
	"os"
)

func main()  {
	client := tui.LoginWindowUI()
	if client == nil {
		os.Exit(1)
	}
	ui := tui.ChatWindowUI(client)
	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}