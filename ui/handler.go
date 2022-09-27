package ui

import (
	"bili/getter"
	"fmt"
	"strings"
	"time"

	"github.com/rivo/tview"
)

func roomInfoHandler(app *tview.Application, roomInfoView *tview.TextView, rankUsersView *tview.TextView, roomInfoChan chan getter.RoomInfo) {
	for roomInfo := range roomInfoChan {
		roomInfoView.SetText(
			roomInfo.Title + "\n" +
				fmt.Sprintf("ID: %d", roomInfo.RoomId) + "\n" +
				fmt.Sprintf("分区: %s/%s", roomInfo.ParentAreaName, roomInfo.AreaName) + "\n" +
				fmt.Sprintf("👁️: %d", roomInfo.Online) + "\n" +
				fmt.Sprintf("❤️: %d", roomInfo.Attention) + "\n" +
				fmt.Sprintf("⏳: %s", roomInfo.Time) + "\n",
		)
		rankUsersView.SetTitle(fmt.Sprintf("Rank(%d)", len(roomInfo.OnlineRankUsers)))

		rankUserStr := ""
		spec := []string{"👑 ", "🥈 ", "🥉 "}
		for idx, rankUser := range roomInfo.OnlineRankUsers {
			if idx < 3 {
				rankUserStr += spec[idx] + rankUser.Name + "\n"
			} else {
				rankUserStr += "   " + rankUser.Name + "\n"
			}
		}
		strings.TrimRight(rankUserStr, "\n")
		rankUsersView.SetText(rankUserStr)
		app.Draw()
	}
}

var lastMsg = getter.DanmuMsg{}
var lastLine = ""

func danmuHandler(app *tview.Application, messages *tview.TextView, busChan chan getter.DanmuMsg) {
	for msg := range busChan {
		if strings.Trim(msg.Content, " ") == "" {
			continue
		}

		viewStr := messages.GetText(false)
		str := ""
		if lastMsg.Type != msg.Type || lastMsg.Author != msg.Author {
			str += fmt.Sprintf("┌ %s %s", time.Now().Format("15:04"), msg.Author) + "\n"
			str += fmt.Sprintf("└ %s", msg.Content) + "\n"
		} else {
			lines := strings.Split(viewStr, "\n")
			lines[len(lines)-2] = strings.Replace(lines[len(lines)-2], "└ ", "│ ", 1)
			viewStr = strings.Join(lines, "\n")
			str += fmt.Sprintf("└ %s", msg.Content) + "\n"
		}
		messages.SetText(viewStr + strings.TrimRight(str, "\n"))
		lastMsg = msg
		app.Draw()
	}
}
