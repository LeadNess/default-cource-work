package tui

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/LeadNess/net-tools/chat/server"
	"github.com/marcusolsson/tui-go"
)

func ServerLogsUI(chatServer *server.TcpChatServer) tui.UI {
	sidebar := tui.NewVBox()
	users := strings.Join(chatServer.ClientsUsernames(), "\n")
	sidebar.Append(tui.NewLabel(users + "\n       "))

	sidebar.SetTitle("Clients")
	sidebar.SetBorder(true)

	history := tui.NewVBox()

	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)
	historyBox.SetTitle("Logs")

	logs := tui.NewVBox(historyBox)
	logs.SetSizePolicy(tui.Expanding, tui.Expanding)

	root := tui.NewHBox(sidebar, logs)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() {
		ui.Quit()
		os.Exit(0)
	})

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
		for clients := range chatServer.Clients() {
			ui.Update(func() {
				sidebar.Remove(0)
				var buf strings.Builder
				for _, client := range clients {
					buf.WriteString(fmt.Sprintf("%s [%s]\n",
						client.Name, client.Conn.RemoteAddr().String()))
				}
				sidebar.Append(tui.NewLabel(buf.String()))
			})
		}
	}()

	return ui
}