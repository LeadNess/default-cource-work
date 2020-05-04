package main

import (
	"default-cource-work/chat/tui"
	"log"
	"os"
)

func main()  {
	client := tui.LoginWindowUI()
	if client == nil {
		os.Exit(0)
	}
	ui := tui.ChatWindowUI(client)
	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}