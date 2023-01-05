// 聊天室主题

package theme3

import (
	"bili/config"
	"bili/getter"
	"bili/sender"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var submitHistory = []string{}
var submitHistoryIndex = 0

func draw(app *tview.Application, roomId int64, busChan chan getter.DanmuMsg, roomInfoChan chan getter.RoomInfo) *tview.Grid {
	grid := tview.NewGrid().SetRows(1, 1, 0, 1, 1).SetBorders(false)

	roomInfoView := tview.NewTextView().SetDynamicColors(true)
	roomInfoView.SetBackgroundColor(tcell.ColorDefault)

	delimiter1 := tview.NewTextView().SetDynamicColors(true) // 分隔符
	delimiter2 := tview.NewTextView().SetDynamicColors(true) // 分隔符
	delimiter1.SetBorder(false).SetBackgroundColor(tcell.ColorDefault)
	delimiter2.SetBorder(false).SetBackgroundColor(tcell.ColorDefault)

	_, _, width, _ := grid.GetRect()
	str := "[" + config.Config.FrameColor + "]"
	for i := 0; i < width; i++ {
		str = str + "—"
	}
	delimiter1.SetText(str)
	delimiter2.SetText(str)

	messagesView := tview.NewTextView().SetDynamicColors(true)
	messagesView.SetBackgroundColor(tcell.ColorDefault)

	input := tview.NewInputField()
	input.SetFormAttributes(0, tcell.ColorDefault, tcell.ColorDefault, tcell.ColorDefault, tcell.ColorDefault)

	grid.
		AddItem(roomInfoView, 0, 0, 1, 1, 0, 0, false).
		AddItem(delimiter1, 1, 0, 1, 1, 0, 0, false).
		AddItem(messagesView, 2, 0, 1, 1, 0, 0, false).
		AddItem(delimiter2, 3, 0, 1, 1, 0, 0, false).
		AddItem(input /*  */, 4, 0, 1, 1, 0, 0, true)

	go roomInfoHandler(app, roomInfoView, roomInfoChan)
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

	grid.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		str := "[" + config.Config.FrameColor + "]"
		for i := 0; i < width; i++ {
			str = str + "—"
		}
		delimiter1.SetText(str)
		delimiter2.SetText(str)
		return x, y, width, height
	})

	return grid
}

func Run(busChan chan getter.DanmuMsg, roomInfoChan chan getter.RoomInfo) {
	app := tview.NewApplication()
	if err := app.SetRoot(draw(app, config.Config.RoomId, busChan, roomInfoChan), true).EnableMouse(false).Run(); err != nil {
		panic(err)
	}
}
