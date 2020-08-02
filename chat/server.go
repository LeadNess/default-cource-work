package main

import (
	"log"
	"os"

	"github.com/LeadNess/net-tools/chat/tui"
)

func main()  {
	server := tui.RunServerUI()
	if server == nil {
		os.Exit(0)
	}
	ui := tui.ServerLogsUI(server)
	go server.Start()
	defer server.Close()
	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
