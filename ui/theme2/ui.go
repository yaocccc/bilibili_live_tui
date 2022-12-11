// 极简主题

package theme2

import (
	"bili/getter"
	"bili/sender"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var submitHistory = []string{}
var submitHistoryIndex = 0

func setBoxAttr(box *tview.Box, title string) {
	box.SetTitleAlign(tview.AlignLeft)
	box.SetBackgroundColor(tcell.ColorDefault)
}

func drawChat() (*tview.Grid, *tview.InputField, *tview.TextView) {
	chatGrid := tview.NewGrid().SetRows(0, 1).SetBorders(false)
	messagesView := tview.NewTextView()
	setBoxAttr(messagesView.Box, "Messages")

	input := tview.NewInputField()
	input.SetFormAttributes(0, tcell.ColorDefault, tcell.ColorDefault, tcell.ColorDefault, tcell.ColorDefault)
	setBoxAttr(input.Box, "Send")

	chatGrid.
		AddItem(messagesView, 0, 0, 1, 1, 0, 0, false).
		AddItem(input, 1, 0, 1, 1, 0, 0, true)

	return chatGrid, input, messagesView
}

func draw(app *tview.Application, roomId int64, busChan chan getter.DanmuMsg, roomInfoChan chan getter.RoomInfo) *tview.Grid {
	chatGrid, input, messagesView := drawChat()
	rootGrid := tview.NewGrid().SetColumns(0).SetBorders(false)
	rootGrid.
		AddItem(chatGrid, 0, 0, 1, 1, 0, 0, true)

	go danmuHandler(app, messagesView, busChan)

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

func Run(roomId int64, busChan chan getter.DanmuMsg, roomInfoChan chan getter.RoomInfo) {
	app := tview.NewApplication()
	if err := app.SetRoot(draw(app, roomId, busChan, roomInfoChan), true).EnableMouse(false).Run(); err != nil {
		panic(err)
	}
}
