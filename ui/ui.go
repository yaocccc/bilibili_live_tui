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
	timeLabel      *tui.Label
}

var submitHistory = []string{}
var submitHistoryIndex = 0

func layoutSidebar(roomInfoChan chan getter.RoomInfo) (tui.Widget, *tui.Box, RoomInfoLabels, *tui.Box) {
	labels := RoomInfoLabels{
		titleLabel:     tui.NewLabel("--------"),
		roomIdLabel:    tui.NewLabel("ID: -----"),
		areaLabel:      tui.NewLabel("----/----"),
		onlineLabel:    tui.NewLabel("👀: --"),
		attentionLabel: tui.NewLabel("❤️ : --"),
		timeLabel:      tui.NewLabel("🕒: --"),
	}

	roomInfo := tui.NewVBox(
		labels.titleLabel,
		labels.roomIdLabel,
		labels.areaLabel,
		tui.NewLabel(""),
		labels.onlineLabel,
		labels.attentionLabel,
		labels.timeLabel,
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

func layoutChat(roomId int64, busChan chan getter.DanmuMsg) (chat *tui.Box, history *tui.Box, input *tui.Entry) {
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

		submitHistory = append(submitHistory, e.Text())
		if len(submitHistory) > 10 {
			submitHistory = submitHistory[1:]
		}
		submitHistoryIndex = len(submitHistory)

		input.SetText("")
	})

	return chat, history, input
}

func roomInfoHandler(ui tui.UI, rankUsersBox *tui.Box, roomInfoLabels RoomInfoLabels, rankUsers *tui.Box, roomInfoChan chan getter.RoomInfo) {
	for roomInfo := range roomInfoChan {
		roomInfoLabels.titleLabel.SetText(roomInfo.Title)
		roomInfoLabels.roomIdLabel.SetText(fmt.Sprintf("ID: %d", roomInfo.RoomId))
		roomInfoLabels.areaLabel.SetText(fmt.Sprintf("%s/%s", roomInfo.ParentAreaName, roomInfo.AreaName))
		roomInfoLabels.onlineLabel.SetText(fmt.Sprintf("👀: %d", roomInfo.Online))
		roomInfoLabels.attentionLabel.SetText(fmt.Sprintf("❤️ : %d", roomInfo.Attention))
		roomInfoLabels.timeLabel.SetText(fmt.Sprintf("🕒: %s", roomInfo.Time))
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

var lastMsg = getter.DanmuMsg{}

func danmuHandler(ui tui.UI, history *tui.Box, lastLabel *tui.Label, roomId int64, busChan chan getter.DanmuMsg) {
	for msg := range busChan {
		if strings.Trim(msg.Content, " ") == "" {
			continue
		}
		// 如果是一条新的用户或新的类型的弹幕，新建一个会话
		if lastMsg.Type != msg.Type || lastMsg.Author != msg.Author {
			history.Append(
				tui.NewLabel(fmt.Sprintf("┌─ %s %s", time.Now().Format("15:04"), msg.Author)),
			)
		} else {
			if lastLabel != nil {
				lastLabel.SetText(strings.Replace(lastLabel.Text(), "└─ ", "│  ", 1))
				lastLabel.SetStyleName("")
			}
		}

		label := tui.NewLabel(fmt.Sprintf("└─ %s", msg.Content))
		history.Append(label)
		lastLabel = label

		lastMsg = msg
		ui.Update(func() {})
	}
}

func Run(roomId int64, busChan chan getter.DanmuMsg, roomInfoChan chan getter.RoomInfo) {
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
	ui.SetKeybinding("up", func() {
		if len(submitHistory) == 0 {
			return
		}
		submitHistoryIndex -= 1
		if submitHistoryIndex < 0 {
			submitHistoryIndex = len(submitHistory) - 1
		}
		input.SetText(submitHistory[submitHistoryIndex])
	})
	ui.SetKeybinding("down", func() {
		if len(submitHistory) == 0 {
			return
		}
		submitHistoryIndex += 1
		if submitHistoryIndex > len(submitHistory)-1 {
			submitHistoryIndex = 0
		}
		input.SetText(submitHistory[submitHistoryIndex])
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
