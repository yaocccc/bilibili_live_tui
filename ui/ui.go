package ui

import (
	"bili/sender"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/marcusolsson/tui-go"
)

func Run(roomId int64, busChan chan []string) {
	sidebar := tui.NewVBox(
		tui.NewLabel(""),
		tui.NewLabel("直播间"),
		tui.NewLabel(fmt.Sprintf("%d", roomId)),
		tui.NewLabel(""),
		tui.NewLabel("coding"),
		tui.NewSpacer(),
	)
	sidebar.SetBorder(true)
	sidebar.SetTitle("Room")

	history := tui.NewVBox()

	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)
	historyBox.SetTitle("History")

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)
	inputBox.SetTitle("Send")

	chat := tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	input.OnSubmit(func(e *tui.Entry) {
		sender.SendMsg(roomId, e.Text())
		input.SetText("")
	})

	root := tui.NewHBox(sidebar, chat)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	history.Append(tui.NewLabel("."))
	var lastLabel *tui.Label
	go func() {
		for msg := range busChan {
			if strings.Trim(msg[1], " ") == "" {
				continue
			}
			if lastLabel != nil {
				lastLabel.SetText(strings.Replace(lastLabel.Text(), "└─ ", "├─ ", 1))
				lastLabel.SetStyleName("")
			}
			label1 := tui.NewLabel(fmt.Sprintf("├─ %s %s", time.Now().Format("15:04"), msg[0]))
			label2 := tui.NewLabel(fmt.Sprintf("└─ %s", msg[1]))
			history.Append(label1)
			history.Append(label2)
			lastLabel = label2
			ui.Update(func() {})
		}
	}()

	ui.SetKeybinding("Esc", func() { ui.Quit() })
	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
