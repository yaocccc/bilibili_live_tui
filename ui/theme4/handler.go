package theme4

import (
	"bili/config"
	"bili/getter"
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

func roomInfoHandler(app *tview.Application, roomInfoView *tview.TextView, rankUsersView *tview.TextView, roomInfoChan chan getter.RoomInfo) {
	for roomInfo := range roomInfoChan {
		roomInfoView.SetText(
			"[" + config.Config.InfoColor + "]" +
				roomInfo.Title + "\n" +
				fmt.Sprintf("ID: %d", roomInfo.RoomId) + "\n" +
				fmt.Sprintf("分区: %s/%s", roomInfo.ParentAreaName, roomInfo.AreaName) + "\n" +
				fmt.Sprintf("👀: %d", roomInfo.Online) + "\n" +
				fmt.Sprintf("❤️: %d", roomInfo.Attention) + "\n" +
				fmt.Sprintf("🕒: %s", roomInfo.Time) + "\n",
		)
		rankUsersView.SetTitle(fmt.Sprintf("Rank(%d)", len(roomInfo.OnlineRankUsers)))

		rankUserStr := ""
		spec := []string{"👑 ", "🥈 ", "🥉 "}
		for idx, rankUser := range roomInfo.OnlineRankUsers {
			rankUserStr += "[" + config.Config.RankColor + "]"
			if idx < 3 {
				rankUserStr += spec[idx] + rankUser.Name + "\n"
			} else {
				rankUserStr += "   " + rankUser.Name + "\n"
			}
		}
		strings.TrimRight(rankUserStr, "\n")
		rankUsersView.SetText(rankUserStr)
		// 滚动到顶部 避免过长显示下半部分
		roomInfoView.ScrollToBeginning()
		rankUsersView.ScrollToBeginning()
		app.Draw()
	}
}

var lastMsg = getter.DanmuMsg{}
var lastLine = ""

func danmuHandler(app *tview.Application, messages *tview.TextView, access *tview.TextView, gift *tview.TextView, busChan chan getter.DanmuMsg) {
	for msg := range busChan {
		if strings.Trim(msg.Content, " ") == "" {
			continue
		}

		viewStr := ""
		str := ""

		// 留意前面的空格显示
		timeStr := msg.Time.Format(" 15:04")
		if config.Config.ShowTime == 0 {
			timeStr = ""
		}

		if config.Config.SingleLine == 1 {
			str += fmt.Sprintf("[%s]%s [%s]%s[%s] %s", config.Config.TimeColor, timeStr, config.Config.NameColor, msg.Author, config.Config.ContentColor, msg.Content)
		} else {
			if lastMsg.Type != msg.Type || lastMsg.Author != msg.Author || lastMsg.Time.Format("15:04") != msg.Time.Format("15:04") {
				str += fmt.Sprintf("[%s]%s [%s]%s[%s]", config.Config.TimeColor, timeStr, config.Config.NameColor, msg.Author, config.Config.ContentColor) + "\n"
			}
			str += fmt.Sprintf(" %s", msg.Content) + "\n"
		}

		switch msg.Type {
		case "INTERACT_WORD":
			viewStr = access.GetText(false)
			access.SetText(viewStr + strings.TrimRight(str, "\n"))
			break
		case "SEND_GIFT":
			viewStr = gift.GetText(false)
			gift.SetText(viewStr + strings.TrimRight(str, "\n"))
			break
		default:
			viewStr = messages.GetText(false)
			messages.SetText(viewStr + strings.TrimRight(str, "\n"))
		}
		lastMsg = msg
		app.Draw()
	}
}
