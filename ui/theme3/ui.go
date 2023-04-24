// simple

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
var bg = tcell.ColorDefault

func draw(app *tview.Application, roomId int64, busChan chan getter.DanmuMsg, roomInfoChan chan getter.RoomInfo) *tview.Grid {
	grid := tview.NewGrid().SetRows(1, 1, 0, 1, 1).SetBorders(false)

	roomInfoView := tview.NewTextView().SetDynamicColors(true)
	roomInfoView.SetBackgroundColor(bg)

	delimiter1 := tview.NewTextView().SetTextColor(tcell.GetColor(config.Config.FrameColor))
	delimiter2 := tview.NewTextView().SetTextColor(tcell.GetColor(config.Config.FrameColor))
	delimiter1.SetBackgroundColor(bg).SetBorder(false)
	delimiter2.SetBackgroundColor(bg).SetBorder(false)

	_, _, width, _ := grid.GetRect()
	str := ""
	for i := 0; i < width; i++ {
		str = str + "—"
	}
	delimiter1.SetText(str)
	delimiter2.SetText(str)

	messagesView := tview.NewTextView().SetDynamicColors(true)
	messagesView.SetBackgroundColor(bg)

	input := tview.NewInputField()
	input.SetFormAttributes(0, tcell.ColorDefault, bg, tcell.ColorDefault, bg)

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
		str := ""
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
	if config.Config.Background != "NONE" {
		bg = tcell.GetColor(config.Config.Background)
	}
	app := tview.NewApplication()
	if err := app.SetRoot(draw(app, config.Config.RoomId, busChan, roomInfoChan), true).EnableMouse(false).Run(); err != nil {
		panic(err)
	}
}
