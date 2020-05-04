package main

import (
	"default-cource-work/chat/tui"
	"log"
	"os"
)

func main()  {
	server := tui.RunServerUI()
	if server == nil {
		os.Exit(0)
	}
	ui := tui.ServerLogsUI(server)
	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
	server.Start()
	defer server.Close()
}