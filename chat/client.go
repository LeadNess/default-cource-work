package main

import (
	"bufio"
	"default-cource-work/chat/client"
	"fmt"
	"github.com/marcusolsson/tui-go"
	"log"
	"os"
	"time"
)

func main()  {
	c := client.NewClient()
	if err := c.Dial(":3333"); err != nil {
		log.Printf("Connect error: %v", err)
	}

	go c.Start()

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter your chat name: ")
	name, _ := reader.ReadString('\n')
	if err := c.SetName(name[:len(name) - 1]); err != nil {
		log.Printf("Set name error: %v", err)
	}


	sidebar := tui.NewVBox()
	for _, user := range <-c.ChatUsers() {
		sidebar.Append(tui.NewLabel(user))
	}
	sidebar.Append(tui.NewLabel("       "))

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
			input.SetText(fmt.Sprintf("Send message error: %v", err))
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

	/*go func() {
		for usersSlice := range c.ChatUsers() {
			ui.Update(func() {
				sidebar.
				for i, user := range usersSlice{
					sidebar.Remove(i)
					sidebar.Insert(i, tui.NewLabel(user))
				}
				sidebar.Remove(len(usersSlice))
				sidebar.Insert(len(usersSlice), tui.NewLabel("       "))
			})
		}
	}()*/

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}