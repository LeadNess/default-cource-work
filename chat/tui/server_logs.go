package tui

import (
	"default-cource-work/chat/server"
	"fmt"
	"github.com/marcusolsson/tui-go"
	"log"
	"strings"
)

func ServerLogsUI(chatServer *server.TcpChatServer) tui.UI {
	sidebar := tui.NewVBox()
	users := strings.Join(chatServer.ClientsUsernames(), "\n")
	sidebar.Append(tui.NewLabel(users + "\n    "))

	sidebar.SetTitle("Clients")
	sidebar.SetBorder(true)

	history := tui.NewVBox()

	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	logs := tui.NewVBox(historyBox)
	logs.SetSizePolicy(tui.Expanding, tui.Expanding)

	root := tui.NewHBox(sidebar, logs)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	go func() {
		for logString := range chatServer.Logs() {
			ui.Update(func() {
				history.Append(tui.NewHBox(
					tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("%s", logString))),
					tui.NewSpacer(),
				))
			})
		}
	}()

	go func() {
		for usersSlice := range chatServer.Clients() {
			ui.Update(func() {
				sidebar.Remove(0)
				users := strings.Join(usersSlice, "\n")
				sidebar.Append(tui.NewLabel(users + "\n    "))
			})
		}
	}()

	return ui
}