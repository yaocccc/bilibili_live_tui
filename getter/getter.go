package getter

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	bg "github.com/iyear/biligo"

	"github.com/asmcos/requests"
	"github.com/tidwall/gjson"
)

type DanmuClient struct {
	roomID        uint32
	auth          bg.CookieAuth
	conn          *websocket.Conn
	unzlibChannel chan []byte
}

type OnlineRankUser struct {
	Name  string
	Score int64
	Rank  int64
}

type RoomInfo struct {
	RoomId          int
	Title           string
	ParentAreaName  string
	AreaName        string
	Online          int64
	Attention       int64
	OnlineRankUsers []OnlineRankUser
}

type receivedInfo struct {
	Cmd        string                 `json:"cmd"`
	Data       map[string]interface{} `json:"data"`
	Info       []interface{}          `json:"info"`
	Full       map[string]interface{} `json:"full"`
	Half       map[string]interface{} `json:"half"`
	Side       map[string]interface{} `json:"side"`
	RoomID     uint32                 `json:"roomid"`
	RealRoomID uint32                 `json:"real_roomid"`
	MsgCommon  string                 `json:"msg_common"`
	MsgSelf    string                 `json:"msg_self"`
	LinkUrl    string                 `json:"link_url"`
	MsgType    string                 `json:"msg_type"`
	ShieldUID  string                 `json:"shield_uid"`
	BusinessID string                 `json:"business_id"`
	Scatter    map[string]interface{} `json:"scatter"`
}

type handShakeInfo struct {
	UID       uint8  `json:"uid"`
	Roomid    uint32 `json:"roomid"`
	Protover  uint8  `json:"protover"`
	Platform  string `json:"platform"`
	Clientver string `json:"clientver"`
	Type      uint8  `json:"type"`
	Key       string `json:"key"`
}

func (d *DanmuClient) connect() {
	var getDanmuInfo = "https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo?id=%d&type=0"
	r, err := requests.Get(fmt.Sprintf(getDanmuInfo, d.roomID))
	if err != nil {
		fmt.Println("request.Get DanmuInfo: ", err)
	}
	fmt.Println("获取弹幕服务器")
	token := gjson.Get(r.Text(), "data.token").String()
	hostList := []string{}
	gjson.Get(r.Text(), "data.host_list").ForEach(func(key, value gjson.Result) bool {
		hostList = append(hostList, value.Get("host").String())
		return true
	})
	hsInfo := handShakeInfo{
		UID:       0,
		Roomid:    d.roomID,
		Protover:  2,
		Platform:  "web",
		Clientver: "1.10.2",
		Type:      2,
		Key:       token,
	}
	for _, h := range hostList {
		d.conn, _, err = websocket.DefaultDialer.Dial(fmt.Sprintf("wss://%s:443/sub", h), nil)
		if err != nil {
			fmt.Println("websocket.Dial: ", err)
			continue
		}
		fmt.Printf("连接弹幕服务器[%s]成功\n", hostList[0])
		break
	}
	if err != nil {
		fmt.Println("websocket.Dial Error")
	}
	jm, err := json.Marshal(hsInfo)
	if err != nil {
		fmt.Println("json.Marshal: ", err)
	}
	err = d.sendPackage(0, 16, 1, 7, 1, jm)
	if err != nil {
		fmt.Println("Conn SendPackage: ", err)
	}
	fmt.Printf("连接房间[%d]成功\n", d.roomID)
}

func (d *DanmuClient) heartBeat() {
	for {
		obj := []byte("5b6f626a656374204f626a6563745d")
		if err := d.sendPackage(0, 16, 1, 2, 1, obj); err != nil {
			fmt.Println("heart beat err: ", err)
			continue
		}
		time.Sleep(30 * time.Second)
	}
}
func (d *DanmuClient) receiveRawMsg(busChan chan []string) {
	for {
		_, rawMsg, _ := d.conn.ReadMessage()
		if rawMsg[7] == 2 {
			msgs := splitMsg(zlibUnCompress(rawMsg[16:]))
			for _, msg := range msgs {
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
				default: // "LIVE" "ACTIVITY_BANNER_UPDATE_V2" "ONLINE_RANK_COUNT" "ONLINE_RANK_TOP3" "ONLINE_RANK_V2" "PANEL" "PREPARING" "WIDGET_BANNER"
					continue
				}
				busChan <- m
			}
		}
	}
}

func (d *DanmuClient) syncRoomInfo(roomInfoChan chan RoomInfo) {
	for {
		roomInfoApi := fmt.Sprintf("https://api.live.bilibili.com/room/v1/Room/get_info?room_id=%d", d.roomID)
		onlineRankApi := fmt.Sprintf("https://api.live.bilibili.com/xlive/general-interface/v1/rank/getOnlineGoldRank?ruid=%s&roomId=%d&page=1&pageSize=50", d.auth.DedeUserID, d.roomID)

		roomInfo := new(RoomInfo)
		roomInfo.OnlineRankUsers = make([]OnlineRankUser, 0)
		r1, err1 := requests.Get(roomInfoApi)
		r2, err2 := requests.Get(onlineRankApi)
		if err1 == nil {
			roomInfo.RoomId = int(d.roomID)
			roomInfo.Title = gjson.Get(r1.Text(), "data.title").String()
			roomInfo.AreaName = gjson.Get(r1.Text(), "data.area_name").String()
			roomInfo.ParentAreaName = gjson.Get(r1.Text(), "data.parent_area_name").String()
			roomInfo.Online = gjson.Get(r1.Text(), "data.online").Int()
			roomInfo.Attention = gjson.Get(r1.Text(), "data.attention").Int()
		}
		if err2 == nil {
			rawUsers := gjson.Get(r2.Text(), "data.OnlineRankItem").Array()
			for _, rawUser := range rawUsers {
				user := OnlineRankUser{
					Name:  rawUser.Get("name").String(),
					Score: rawUser.Get("score").Int(),
					Rank:  rawUser.Get("userRank").Int(),
				}
				roomInfo.OnlineRankUsers = append(roomInfo.OnlineRankUsers, user)
			}
		}

		roomInfoChan <- *roomInfo
		time.Sleep(30 * time.Second)
	}
}

func Run(roomID int64, auth bg.CookieAuth, busChan chan []string, roomInfoChan chan RoomInfo) {
	dc := DanmuClient{
		roomID:        uint32(roomID),
		auth:          auth,
		conn:          new(websocket.Conn),
		unzlibChannel: make(chan []byte, 100),
	}
	dc.connect()
	go dc.heartBeat()
	go dc.receiveRawMsg(busChan)
	go dc.syncRoomInfo(roomInfoChan)
}
