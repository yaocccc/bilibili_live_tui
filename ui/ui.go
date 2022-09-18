package ui

import (
	"bili/getter"
	"bili/sender"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/marcusolsson/tui-go"
)

type RoomInfoLabels struct {
	titleLabel     *tui.Label
	roomIdLabel    *tui.Label
	areaLabel      *tui.Label
	onlineLabel    *tui.Label
	attentionLabel *tui.Label
}

func layoutSidebar(roomInfoChan chan getter.RoomInfo) (tui.Widget, *tui.Box, RoomInfoLabels, *tui.Box) {
	labels := RoomInfoLabels{
		titleLabel:     tui.NewLabel("--------"),
		roomIdLabel:    tui.NewLabel("ID: -----"),
		areaLabel:      tui.NewLabel("----/----"),
		onlineLabel:    tui.NewLabel("👀: --"),
		attentionLabel: tui.NewLabel("❤️ : --"),
	}

	roomInfo := tui.NewVBox(
		tui.NewLabel(""),
		labels.titleLabel,
		labels.roomIdLabel,
		labels.areaLabel,
		tui.NewLabel(""),
		labels.onlineLabel,
		labels.attentionLabel,
		tui.NewSpacer(),
	)
	roomInfo.SetBorder(true)
	roomInfo.SetTitle("Room")

	rankUsers := tui.NewVBox()
	rankUsersScroll := tui.NewScrollArea(rankUsers)
	rankUsersScroll.SetAutoscrollToBottom(false)

	rankUsersBox := tui.NewVBox(rankUsersScroll)
	rankUsersBox.SetBorder(true)
	rankUsersBox.SetTitle("Rank(0)")

	sidebar := tui.NewVBox(roomInfo, rankUsersBox)
	return sidebar, rankUsersBox, labels, rankUsers
}

func layoutChat(roomId int64, busChan chan []string) (chat *tui.Box, history *tui.Box, input *tui.Entry) {
	history = tui.NewVBox()

	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)
	historyBox.SetTitle("History")

	input = tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)
	inputBox.SetTitle("Send")

	chat = tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	input.OnSubmit(func(e *tui.Entry) {
		go sender.SendMsg(roomId, e.Text(), busChan)
		input.SetText("")
	})

	history.Append(tui.NewLabel("."))

	return chat, history, input
}

func roomInfoHandler(ui tui.UI, rankUsersBox *tui.Box, roomInfoLabels RoomInfoLabels, rankUsers *tui.Box, roomInfoChan chan getter.RoomInfo) {
	for roomInfo := range roomInfoChan {
		roomInfoLabels.titleLabel.SetText(roomInfo.Title)
		roomInfoLabels.roomIdLabel.SetText(fmt.Sprintf("ID: %d", roomInfo.RoomId))
		roomInfoLabels.areaLabel.SetText(fmt.Sprintf("%s/%s", roomInfo.ParentAreaName, roomInfo.AreaName))
		roomInfoLabels.onlineLabel.SetText(fmt.Sprintf("👀: %d", roomInfo.Online))
		roomInfoLabels.attentionLabel.SetText(fmt.Sprintf("❤️ : %d", roomInfo.Attention))
		rankUsersBox.SetTitle(fmt.Sprintf("Rank(%d)", len(roomInfo.OnlineRankUsers)))

		for rankUsers.Length() > 0 {
			rankUsers.Remove(0)
		}
		spec := []string{"👑 ", "🥈 ", "🥉 "}
		for idx, rankUser := range roomInfo.OnlineRankUsers {
			if idx < 3 {
				rankUsers.Append(tui.NewLabel(spec[idx] + rankUser.Name))
			} else {
				rankUsers.Append(tui.NewLabel("   " + rankUser.Name))
			}
		}
		ui.Update(func() {})
	}
}

func danmuHandler(ui tui.UI, history *tui.Box, lastLabel *tui.Label, roomId int64, busChan chan []string) {
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
}

func Run(roomId int64, busChan chan []string, roomInfoChan chan getter.RoomInfo) {
	sidebar, rankUsersBox, roomInfoLabels, rankUsers := layoutSidebar(roomInfoChan)
	chat, history, input := layoutChat(roomId, busChan)

	root := tui.NewHBox(sidebar, chat)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	var lastLabel *tui.Label
	go danmuHandler(ui, history, lastLabel, roomId, busChan)
	go roomInfoHandler(ui, rankUsersBox, roomInfoLabels, rankUsers, roomInfoChan)

	ui.SetKeybinding("Esc", func() { ui.Quit() })
	ui.SetKeybinding("Ctrl+c", func() { ui.Quit() })
	ui.SetKeybinding("Ctrl+u", func() { input.SetText("") })

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
