// 聊天室主题

package theme4

import (
	"bili/config"
	"bili/getter"
	"bili/sender"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var submitHistory = []string{}
var submitHistoryIndex = 0

func setBoxAttr(box *tview.Box, title string) {
	box.SetBorder(true)
	box.SetTitleAlign(tview.AlignLeft)
	box.SetTitle(title)
	box.SetBackgroundColor(tcell.ColorDefault)
	box.SetBorderColor(tcell.GetColor(config.Config.FrameColor))
	box.SetTitleColor(tcell.GetColor(config.Config.FrameColor))
}

func drawSlidebar() (*tview.Grid, *tview.TextView, *tview.TextView) {
	slidebarGrid := tview.NewGrid().SetRows(0, 0).SetBorders(false)
	roomInfoView := tview.NewTextView().SetDynamicColors(true)
	setBoxAttr(roomInfoView.Box, "RoomInfo")

	rankUsersView := tview.NewTextView().SetDynamicColors(true)
	setBoxAttr(rankUsersView.Box, "RankUsers")

	slidebarGrid.
		AddItem(roomInfoView, 0, 0, 1, 1, 0, 0, false).
		AddItem(rankUsersView, 1, 0, 1, 1, 0, 0, false)

	return slidebarGrid, roomInfoView, rankUsersView
}

func drawChat() (*tview.Grid, *tview.InputField, *tview.TextView, *tview.TextView, *tview.TextView) {
	chatGrid := tview.NewGrid().SetRows(0, 3).SetBorders(false)
	messagesView := tview.NewTextView().SetDynamicColors(true)
	setBoxAttr(messagesView.Box, "Messages")
	accessView := tview.NewTextView().SetDynamicColors(true)
	setBoxAttr(accessView.Box, "Access")
	giftView := tview.NewTextView().SetDynamicColors(true)
	setBoxAttr(giftView.Box, "Gift")

	input := tview.NewInputField()
	input.SetFormAttributes(0, tcell.ColorDefault, tcell.ColorDefault, tcell.ColorDefault, tcell.ColorDefault)
	setBoxAttr(input.Box, "Send")

	chatGrid.
		AddItem(messagesView, 0, 0, 1, 1, 0, 0, false).
		AddItem(accessView, 0, 1, 1, 1, 0, 0, false).
		AddItem(giftView, 0, 2, 1, 1, 0, 0, false).
		AddItem(input, 1, 0, 1, 3, 0, 0, true)

	return chatGrid, input, messagesView, accessView, giftView
}

func draw(app *tview.Application, roomId int64, busChan chan getter.DanmuMsg, roomInfoChan chan getter.RoomInfo) *tview.Grid {
	slidebarGrid, roomInfoView, rankUsersView := drawSlidebar()
	chatGrid, input, messagesView, accessView, giftView := drawChat()
	rootGrid := tview.NewGrid().SetColumns(20, 0).SetBorders(false)
	rootGrid.
		AddItem(slidebarGrid, 0, 0, 1, 1, 0, 0, false).
		AddItem(chatGrid, 0, 1, 1, 1, 0, 0, true)

	go roomInfoHandler(app, roomInfoView, rankUsersView, roomInfoChan)
	go danmuHandler(app, messagesView, accessView, giftView, busChan)

	input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			go sender.SendMsg(roomId, input.GetText(), busChan)

			submitHistory = append(submitHistory, input.GetText())
			if len(submitHistory) > 10 {
				submitHistory = submitHistory[1:]
			}
			submitHistoryIndex = len(submitHistory)

			input.SetText("")
		}
	})

	return rootGrid
}

func Run(busChan chan getter.DanmuMsg, roomInfoChan chan getter.RoomInfo) {
	app := tview.NewApplication()
	if err := app.SetRoot(draw(app, config.Config.RoomId, busChan, roomInfoChan), true).EnableMouse(false).Run(); err != nil {
		panic(err)
	}
}
