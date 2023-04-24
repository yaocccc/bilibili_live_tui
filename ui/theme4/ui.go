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
var bg = tcell.ColorDefault

func setBoxAttr(box *tview.Box, title string) {
	box.SetBorder(true)
	box.SetTitleAlign(tview.AlignLeft)
	box.SetTitle(title)
	box.SetBackgroundColor(bg)
	box.SetBorderColor(tcell.GetColor(config.Config.FrameColor))
	box.SetTitleColor(tcell.GetColor(config.Config.FrameColor))
}

func drawSlidebar() (*tview.Grid, *tview.TextView, *tview.TextView) {
	slidebarGrid := tview.NewGrid().SetRows(0, 0).SetBorders(false)
	roomInfoView := tview.NewTextView().SetDynamicColors(true)
	roomInfoView.SetBackgroundColor(bg)
	setBoxAttr(roomInfoView.Box, "RoomInfo")

	rankUsersView := tview.NewTextView().SetDynamicColors(true)
	rankUsersView.SetBackgroundColor(bg)
	setBoxAttr(rankUsersView.Box, "RankUsers")

	slidebarGrid.
		AddItem(roomInfoView, 0, 0, 1, 1, 0, 0, false).
		AddItem(rankUsersView, 1, 0, 1, 1, 0, 0, false)

	return slidebarGrid, roomInfoView, rankUsersView
}

func drawChat() (*tview.Grid, *tview.InputField, *tview.TextView, *tview.TextView, *tview.TextView) {
	chatGrid := tview.NewGrid().SetRows(0, 0, 3).SetBorders(false)
	messagesView := tview.NewTextView().SetDynamicColors(true)
	messagesView.SetBackgroundColor(bg)
	setBoxAttr(messagesView.Box, "Messages")
	accessView := tview.NewTextView().SetDynamicColors(true)
	accessView.SetBackgroundColor(bg)
	setBoxAttr(accessView.Box, "Access")
	giftView := tview.NewTextView().SetDynamicColors(true)
	giftView.SetBackgroundColor(bg)
	setBoxAttr(giftView.Box, "Gift")

	input := tview.NewInputField()
	input.SetFormAttributes(0, tcell.ColorDefault, bg, tcell.ColorDefault, bg)
	setBoxAttr(input.Box, "Send")

	// 动态布局 宽度大于80时 采用三列布局 否则采用两列布局
	// 三列 | 消息 | 访问 | 礼物 |, 两列 | 消息 | 访问 / 礼物 |
	chatGrid.
		AddItem(messagesView, 0, 0, 2, 1, 0, 0, false).  // 小于80时 | danmu | access / gift |
		AddItem(messagesView, 0, 0, 2, 1, 0, 80, false). // 超过80时 | danmu | access | gift |
		AddItem(accessView, 0, 1, 1, 1, 0, 0, false).    // 小于80时 | danmu | access / gift |
		AddItem(accessView, 0, 1, 2, 1, 0, 80, false).   // 超过80时 | danmu | access | gift |
		AddItem(giftView, 1, 1, 1, 1, 0, 0, false).      // 小于80时 | danmu | access / gift |
		AddItem(giftView, 0, 2, 2, 1, 0, 80, false).     // 超过80时 | danmu | access | gift |
		AddItem(input, 2, 0, 1, 2, 0, 0, true).          // 小于80时 | danmu | access / gift |
		AddItem(input, 2, 0, 1, 3, 0, 80, true)          // 超过80时 | danmu | access | gift |

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
	if config.Config.Background != "NONE" {
		bg = tcell.GetColor(config.Config.Background)
	}
	app := tview.NewApplication()
	if err := app.SetRoot(draw(app, config.Config.RoomId, busChan, roomInfoChan), true).EnableMouse(false).Run(); err != nil {
		panic(err)
	}
}
