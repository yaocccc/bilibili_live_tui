package getter

import (
	"encoding/json"
	"fmt"
)

func (d *DanmuClient) process(busChan chan string) {
	for msg := range d.unzlibChannel {
		uz := msg[16:]
		js := new(receivedInfo)
		json.Unmarshal(uz, js)
		m := ""
		switch js.Cmd {
		case "COMBO_SEND":
			m = fmt.Sprintf("%s 送给 %s %d 个 %s", js.Data["uname"].(string), js.Data["r_uname"].(string), int(js.Data["combo_num"].(float64)), js.Data["gift_name"].(string))
		case "DANMU_MSG":
			m = fmt.Sprintf("%s: %s", js.Info[2].([]interface{})[1].(string), js.Info[1].(string))
		case "GUARD_BUY":
			m = fmt.Sprintf("%s 购买了 %s", js.Data["username"].(string), js.Data["gift_name"].(string))
		case "INTERACT_WORD":
			m = fmt.Sprintf("%s 进入了房间", js.Data["uname"].(string))
		case "SEND_GIFT":
			m = fmt.Sprintf("%s 投喂了 %d 个 %s", js.Data["uname"].(string), int(js.Data["num"].(float64)), js.Data["giftName"].(string))
		case "USER_TOAST_MSG":
			m = js.Data["toast_msg"].(string)
		case "LIVE_INTERACTIVE_GAME":
			continue
		case "NOTICE_MSG":
			m = js.MsgSelf
		case "LIVE":
			continue
		case "ACTIVITY_BANNER_UPDATE_V2":
			continue
		case "ONLINE_RANK_COUNT":
			continue
		case "ONLINE_RANK_TOP3":
			continue
		case "ONLINE_RANK_V2":
			continue
		case "PANEL":
			continue
		case "PREPARING":
			continue
		case "WIDGET_BANNER":
			continue
		}
		busChan <- m
	}
}
