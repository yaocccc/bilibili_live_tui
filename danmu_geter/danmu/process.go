package danmu

import (
	"encoding/json"
	"fmt"
	"time"
)

func (d *DanmuClient) process() {
	for {
		select {
		case msg := <-d.unzlibChannel:
			uz := msg[16:]
			js := new(receivedInfo)
			json.Unmarshal(uz, js)
			switch js.Cmd {
			case "ACTIVITY_BANNER_UPDATE_V2":
				continue
			case "COMBO_SEND":
				fmt.Printf("%s 送给 %s %d 个 %s\n", js.Data["uname"].(string), js.Data["r_uname"].(string), int(js.Data["combo_num"].(float64)), js.Data["gift_name"].(string))
			case "DANMU_MSG":
				fmt.Printf("%s: %s\n", js.Info[2].([]interface{})[1].(string), js.Info[1].(string))
			case "ENTRY_EFFECT":
				fmt.Printf("ENTRY_EFFECT %s\n", js.Data["copy_writing_v2"].(string))
			case "GUARD_BUY":
				fmt.Printf("%s 购买了 %s\n", js.Data["username"].(string), js.Data["gift_name"].(string))
			case "INTERACT_WORD":
				fmt.Printf("%s 进入了房间\n", js.Data["uname"].(string))
			case "LIVE_INTERACTIVE_GAME":
				continue
			case "LIVE":
				continue
			case "NOTICE_MSG":
				fmt.Printf("%s\n", js.MsgSelf)
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
			case "SEND_GIFT":
				fmt.Printf("%s 投喂了 %d 个 %s\n", js.Data["uname"].(string), int(js.Data["num"].(float64)), js.Data["giftName"].(string))
			case "USER_TOAST_MSG":
				fmt.Printf("%s\n", js.Data["toast_msg"].(string))
			case "WIDGET_BANNER":
				continue
			}
		// case msg := <-d.heartBeatChannel:
		// 	fmt.Printf("HeartBeat...实时人气: %d. \n", ByteArrToDecimal(msg[16:]))
		case msg := <-d.serverNoticeChannel:
			if msg[7] == 0 {
				uz := msg[16:]
				js := new(receivedInfo)
				json.Unmarshal(uz, js)
				switch js.Cmd {
				case "NOTICE_MSG":
					fmt.Printf("[%s] From Server %s: %s. \n", time.Now().Format("2006-01-02 15:04:05"), js.Cmd, js.MsgSelf)
				case "ROOM_RANK":
					fmt.Printf("[%s] From Server %s: %d %s. \n", time.Now().Format("2006-01-02 15:04:05"), js.Cmd, uint32(js.Data["roomid"].(float64)), js.Data["rank_desc"].(string))
				case "ROOM_REAL_TIME_MESSAGE_UPDATE":
					fmt.Printf("[%s] From Server %s: [%d] 关注: %d, 粉丝团: %d. \n", time.Now().Format("2006-01-02 15:04:05"), js.Cmd, uint32(js.Data["roomid"].(float64)), int(js.Data["fans"].(float64)), int(js.Data["fans_club"].(float64)))
				case "HOT_RANK_CHANGED":
					fmt.Printf("[%s] From Server %s: [%d] Rank: %d, Trend: %d, Area name: %d, Countdown: %d. \n", time.Now().Format("2006-01-02 15:04:05"), js.Cmd, d.roomID, int(js.Data["rank"].(float64)), int(js.Data["trend"].(float64)), js.Data["area_name"].(string), int(js.Data["countdown"].(float64)))
				}
			} else {
				fmt.Println(msg)
			}
		}
	}
}
