package tui

import (
	"default-cource-work/chat/client"
	"fmt"
	"github.com/marcusolsson/tui-go"
	"log"
)

var clientLogo = `
  ████████████████████ ███████████       █████████ █████              █████         █████████ ████ ███                   █████   
░█░░░███░░░███░░░░░██░░███░░░░░███     ███░░░░░██░░███              ░░███         ███░░░░░██░░███░░░                   ░░███    
░   ░███  ███     ░░░ ░███    ░███    ███     ░░░ ░███████   ██████ ███████      ███     ░░░ ░███████  ██████ ████████ ███████  
    ░███ ░███         ░██████████    ░███         ░███░░███ ░░░░░██░░░███░      ░███         ░██░░███ ███░░██░░███░░██░░░███░   
    ░███ ░███         ░███░░░░░░     ░███         ░███ ░███  ███████ ░███       ░███         ░███░███░███████ ░███ ░███ ░███    
    ░███ ░░███     ███░███           ░░███     ███░███ ░███ ███░░███ ░███ ███   ░░███     ███░███░███░███░░░  ░███ ░███ ░███ ███
    █████ ░░█████████ █████           ░░█████████ ████ ████░░████████░░█████     ░░█████████ ████████░░██████ ████ █████░░█████ 
   ░░░░░   ░░░░░░░░░ ░░░░░             ░░░░░░░░░ ░░░░ ░░░░░ ░░░░░░░░  ░░░░░       ░░░░░░░░░ ░░░░░░░░░ ░░░░░░ ░░░░ ░░░░░  ░░░░░   
`

func LoginWindowUI() *client.TcpChatClient {
	username := tui.NewEntry()
	username.SetFocused(true)

	address := tui.NewEntry()

	form := tui.NewGrid(0, 0)
	form.AppendRow(tui.NewLabel("User"), tui.NewLabel("Server address"))
	form.AppendRow(username, address)

	connect := tui.NewButton("[Connect]")

	button := tui.NewHBox(
		tui.NewSpacer(),
		tui.NewPadder(1, 0, connect),
	)

	info := tui.NewLabel("")

	window := tui.NewVBox(
		tui.NewPadder(10, 1, tui.NewLabel(clientLogo)),
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

	tui.DefaultFocusChain.Set(username, address, connect)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })
	chatClient := client.NewClient()
	connect.OnActivated(func(b *tui.Button) {
		if err := chatClient.Dial(address.Text()); err != nil {
			info.SetText(fmt.Sprintf("Connect error: %v", err))
			return
		}

		go chatClient.Start()

		if err := chatClient.SetName(username.Text()); err != nil {
			info.SetText(fmt.Sprintf("Set name error: %v", err))
			return
		}
		ui.Quit()
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
	return chatClient
}