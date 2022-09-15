package getter

import (
	"encoding/json"
	"fmt"
)

func (d *DanmuClient) process(busChan chan []string) {
	for msg := range d.unzlibChannel {
		uz := msg[16:]
		js := new(receivedInfo)
		json.Unmarshal(uz, js)
		m := make([]string, 2)
		switch js.Cmd {
		case "COMBO_SEND":
			m[0] = js.Data["uname"].(string)
			m[1] = fmt.Sprintf("送给 %s %d 个 %s", js.Data["r_uname"].(string), int(js.Data["combo_num"].(float64)), js.Data["gift_name"].(string))
		case "DANMU_MSG":
			m[0] = js.Info[2].([]interface{})[1].(string)
			m[1] = js.Info[1].(string)
		case "GUARD_BUY":
			m[0] = js.Data["username"].(string)
			m[1] = fmt.Sprintf("购买了 %s", js.Data["giftName"].(string))
		case "INTERACT_WORD":
			m[0] = js.Data["uname"].(string)
			m[1] = "进入了房间"
		case "SEND_GIFT":
			m[0] = js.Data["uname"].(string)
			m[1] = fmt.Sprintf("投喂了 %d 个 %s", int(js.Data["num"].(float64)), js.Data["giftName"].(string))
		case "USER_TOAST_MSG":
			m[0] = "system"
			m[1] = js.Data["toast_msg"].(string)
		case "LIVE_INTERACTIVE_GAME":
			continue
		case "NOTICE_MSG":
			m[0] = "system"
			m[1] = js.MsgSelf
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
