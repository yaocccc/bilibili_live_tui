package theme2

import (
	"bili/config"
	"bili/getter"
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

var lastMsg = getter.DanmuMsg{}
var lastLine = ""

func danmuHandler(app *tview.Application, messages *tview.TextView, busChan chan getter.DanmuMsg) {
	for msg := range busChan {
		if strings.Trim(msg.Content, " ") == "" {
			continue
		}

		viewStr := messages.GetText(false)
		str := ""
		if lastMsg.Type != msg.Type || lastMsg.Author != msg.Author || lastMsg.Time.Format("15:04") != msg.Time.Format("15:04") {
			str += fmt.Sprintf(" [%s]%s [%s]%s[%s]", config.Config.TimeColor, msg.Time.Format("15:04"), config.Config.NameColor, msg.Author, config.Config.ContentColor) + "\n"
		}
		str += fmt.Sprintf(" %s", msg.Content) + "\n"
		messages.SetText(viewStr + strings.TrimRight(str, "\n"))
		lastMsg = msg
		app.Draw()
	}
}
