package ui

import (
	"bili/sender"
	"log"
	"time"

	"github.com/marcusolsson/tui-go"
)

func Run(roomId int64, busChan chan string) {
	sidebar := tui.NewVBox(
		tui.NewLabel("直播间: 123123"),
		tui.NewLabel("coding"),
		tui.NewLabel(""),
		tui.NewLabel("👀 100"),
		tui.NewLabel("🔥 2w"),
		tui.NewSpacer(),
	)
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
		sender.SendMsg(roomId, e.Text())
		input.SetText("")
	})

	root := tui.NewHBox(sidebar, chat)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	go func() {
		for msg := range busChan {
			history.Append(tui.NewHBox(
				tui.NewLabel(time.Now().Format("15:04")),
				tui.NewPadder(1, 0, tui.NewLabel(msg)),
				tui.NewSpacer(),
			))
			ui.Update(func() {})
		}
	}()

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
