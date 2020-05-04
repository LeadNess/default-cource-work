package tui

import (
	"default-cource-work/chat/client"
	"fmt"
	"github.com/marcusolsson/tui-go"
	"log"
	"strings"
	"time"
)

func ChatWindowUI(c *client.TcpChatClient) tui.UI {
	sidebar := tui.NewVBox()
	users := strings.Join(<-c.ChatUsers(), "\n")
	sidebar.Append(tui.NewLabel(users + "\n    "))

	sidebar.SetTitle("Users")
	sidebar.SetBorder(true)

	history := tui.NewVBox()

	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	input.OnSubmit(func(e *tui.Entry) {
		if err := c.SendMessage(e.Text()); err != nil {
			history.Append(tui.NewHBox(
				tui.NewLabel(time.Now().Format("15:04")),
				tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("Send message error: %v", err))),
				tui.NewSpacer(),
			))
		} else {
			input.SetText("")
		}
	})

	root := tui.NewHBox(sidebar, chat)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	go func() {
		for message := range c.Incoming() {
			ui.Update(func() {
				history.Append(tui.NewHBox(
					tui.NewLabel(time.Now().Format("15:04")),
					tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", message.Name))),
					tui.NewLabel(message.Message),
					tui.NewSpacer(),
				))
			})
		}
	}()

	go func() {
		for usersSlice := range c.ChatUsers() {
			ui.Update(func() {
				sidebar.Remove(0)
				users := strings.Join(usersSlice, "\n")
				sidebar.Append(tui.NewLabel(users + "\n    "))
			})
		}
	}()

	return ui
}