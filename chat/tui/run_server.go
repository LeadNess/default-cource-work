package tui

import (
	"fmt"
	"log"
	"os"

	"github.com/LeadNess/net-tools/chat/server"
	"github.com/marcusolsson/tui-go"
)

var serverLogo = `
 ████████████████████ ███████████       █████████ █████              █████        █████████                                             
░█░░░███░░░███░░░░░██░░███░░░░░███     ███░░░░░██░░███              ░░███        ███░░░░░███                                            
░   ░███  ███     ░░░ ░███    ░███    ███     ░░░ ░███████   ██████ ███████     ░███    ░░░  ██████ ████████ █████ ███████████ ████████ 
    ░███ ░███         ░██████████    ░███         ░███░░███ ░░░░░██░░░███░      ░░█████████ ███░░██░░███░░██░░███ ░░██████░░██░░███░░███
    ░███ ░███         ░███░░░░░░     ░███         ░███ ░███  ███████ ░███        ░░░░░░░░██░███████ ░███ ░░░ ░███  ░██░███████ ░███ ░░░ 
    ░███ ░░███     ███░███           ░░███     ███░███ ░███ ███░░███ ░███ ███    ███    ░██░███░░░  ░███     ░░███ ███░███░░░  ░███     
    █████ ░░█████████ █████           ░░█████████ ████ ████░░████████░░█████    ░░█████████░░██████ █████     ░░█████ ░░██████ █████    
   ░░░░░   ░░░░░░░░░ ░░░░░             ░░░░░░░░░ ░░░░ ░░░░░ ░░░░░░░░  ░░░░░      ░░░░░░░░░  ░░░░░░ ░░░░░       ░░░░░   ░░░░░░ ░░░░░     
`

func RunServerUI() *server.TcpChatServer {
	address := tui.NewEntry()
	address.SetFocused(true)

	form := tui.NewGrid(0, 0)
	form.AppendRow(tui.NewLabel("Server address"))
	form.AppendRow(address)

	runServer := tui.NewButton("[Run server]")

	button := tui.NewHBox(
		tui.NewSpacer(),
		tui.NewPadder(1, 0, runServer),
	)

	info := tui.NewLabel("")

	window := tui.NewVBox(
		tui.NewPadder(10, 1, tui.NewLabel(serverLogo)),
		tui.NewPadder(1, 0, info),
		tui.NewPadder(1, 1, form),
		button,
	)
	window.SetBorder(true)

	wrapper := tui.NewVBox(
		tui.NewSpacer(),
		window,
		tui.NewSpacer(),
	)
	content := tui.NewHBox(tui.NewSpacer(), wrapper, tui.NewSpacer())

	root := tui.NewVBox(content)

	tui.DefaultFocusChain.Set(address, runServer)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() {
		ui.Quit()
		os.Exit(0)
	})

	chatServer := server.NewServer()

	runServer.OnActivated(func(b *tui.Button) {
		if err = chatServer.Listen(address.Text()); err != nil {
			info.SetText(fmt.Sprintf("Running server error: %v", err))
			return
		}
		ui.Quit()
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
	return chatServer
}