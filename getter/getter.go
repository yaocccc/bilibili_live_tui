package getter

import (
	"bili/config"
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
	isClosed      bool
}

type OnlineRankUser struct {
	Name  string
	Score int64
	Rank  int64
}

type RoomInfo struct {
	RoomId          int
	Uid             int
	Title           string
	ParentAreaName  string
	AreaName        string
	Online          int64
	Attention       int64
	Time            string
	OnlineRankUsers []OnlineRankUser
}

type DanmuMsg struct {
	Author  string
	Content string
	Type    string
	Time    time.Time
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

func (d *DanmuClient) connect() error {
	var (
		r   *requests.Response
		err error
		jm  []byte
	)

	var getDanmuInfo = "https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo?id=%d&type=0"

	r, err = requests.Get(fmt.Sprintf(getDanmuInfo, d.roomID))
	if err != nil {
		time.Sleep(1 * time.Second)
	}

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
			continue
		}
		break
	}
	if err != nil {
		return err
	}
	jm, err = json.Marshal(hsInfo)
	if err != nil {
		return err
	}

	err = d.sendPackage(0, 16, 1, 7, 1, jm)
	return err
}

func (d *DanmuClient) getHistory(busChan chan DanmuMsg) {
	historyApi := fmt.Sprintf("https://api.live.bilibili.com/xlive/web-room/v1/dM/gethistory?roomid=%d", d.roomID)
	r, err := requests.Get(historyApi)
	if err != nil {
		return
	}

	histories := gjson.Get(r.Text(), "data.room").Array()
	for _, history := range histories {
		t, _ := time.Parse("2006-01-02 15:04:05", history.Get("timeline").String())
		danmu := DanmuMsg{
			Author:  history.Get("nickname").String(),
			Content: history.Get("text").String(),
			Type:    "DANMU_MSG",
			Time:    t,
		}
		busChan <- danmu
	}
}

func (d *DanmuClient) heartBeat(msgChan chan DanmuMsg) {
	for {
		if d.isClosed {
			return
		}
		obj := []byte("5b6f626a656374204f626a6563745d")
		if err := d.sendPackage(0, 16, 1, 2, 1, obj); err != nil {
			msgChan <- DanmuMsg{
				// Author
				// Content string
				// Type    string
			}
			continue
		}
		time.Sleep(30 * time.Second)
	}
}

// fuck 每次启动总容易失败panic
func (d *DanmuClient) receiveRawMsg(busChan chan DanmuMsg) {
	for {
		if d.isClosed {
			return
		}
		_, rawMsg, err := d.conn.ReadMessage()
		if err != nil {
			d.isClosed = true
		}
		if len(rawMsg) >= 8 && rawMsg[7] == 2 {
			msgs := splitMsg(zlibUnCompress(rawMsg[16:]))
			for _, msg := range msgs {
				uz := msg[16:]
				js := new(receivedInfo)
				json.Unmarshal(uz, js)
				m := DanmuMsg{}
				switch js.Cmd {
				case "COMBO_SEND":
					m.Author = js.Data["uname"].(string)
					m.Content = fmt.Sprintf("送给 %s %d 个 %s", js.Data["r_uname"].(string), int(js.Data["combo_num"].(float64)), js.Data["gift_name"].(string))
				case "DANMU_MSG":
					m.Author = js.Info[2].([]interface{})[1].(string)
					m.Content = js.Info[1].(string)
				case "GUARD_BUY":
					m.Author = js.Data["username"].(string)
					m.Content = fmt.Sprintf("购买了 %s", js.Data["giftName"].(string))
				case "INTERACT_WORD":
					m.Author = js.Data["uname"].(string)
					m.Content = "进入了房间"
				case "SEND_GIFT":
					m.Author = js.Data["uname"].(string)
					m.Content = fmt.Sprintf("投喂了 %d 个 %s", int(js.Data["num"].(float64)), js.Data["giftName"].(string))
				case "USER_TOAST_MSG":
					m.Author = "system"
					m.Content = js.Data["toast_msg"].(string)
				case "NOTICE_MSG":
					m.Author = "system"
					m.Content = js.MsgSelf
				default: // "LIVE" "ACTIVITY_BANNER_UPDATE_V2" "ONLINE_RANK_COUNT" "ONLINE_RANK_TOP3" "ONLINE_RANK_V2" "PANEL" "PREPARING" "WIDGET_BANNER" "LIVE_INTERACTIVE_GAME"
					continue
				}
				m.Type = js.Cmd
				m.Time = time.Now()
				busChan <- m
			}
		}
	}
}

func (d *DanmuClient) syncRoomInfo(roomInfoChan chan RoomInfo) {
	for {
		if d.isClosed {
			return
		}

		roomInfoApi := fmt.Sprintf("https://api.live.bilibili.com/room/v1/room/get_info?room_id=%d", d.roomID)
		roomInfo := new(RoomInfo)
		roomInfo.OnlineRankUsers = make([]OnlineRankUser, 0)
		r1, err1 := requests.Get(roomInfoApi)
		if err1 == nil {
			roomInfo.RoomId = int(d.roomID)
			roomInfo.Uid = int(gjson.Get(r1.Text(), "data.uid").Int())
			roomInfo.Title = gjson.Get(r1.Text(), "data.title").String()
			roomInfo.AreaName = gjson.Get(r1.Text(), "data.area_name").String()
			roomInfo.ParentAreaName = gjson.Get(r1.Text(), "data.parent_area_name").String()
			roomInfo.Online = gjson.Get(r1.Text(), "data.online").Int()
			roomInfo.Attention = gjson.Get(r1.Text(), "data.attention").Int()
			_time, _ := time.Parse("2006-01-02 15:04:05", gjson.Get(r1.Text(), "data.live_time").String())
			seconds := time.Now().Unix() - _time.Unix() + 8*60*60
			days := seconds / 86400
			hours := (seconds % 86400) / 3600
			minutes := (seconds % 3600) / 60
			if days > 0 {
				roomInfo.Time = fmt.Sprintf("%d天%d时%d分", days, hours, minutes)
			} else if hours > 0 {
				roomInfo.Time = fmt.Sprintf("%d时%d分", hours, minutes)
			} else {
				roomInfo.Time = fmt.Sprintf("%d分", minutes)
			}
		}

		onlineRankApi := fmt.Sprintf("https://api.live.bilibili.com/xlive/general-interface/v1/rank/getOnlineGoldRank?ruid=%d&roomId=%d&page=1&pageSize=50", roomInfo.Uid, d.roomID)
		r2, err2 := requests.Get(onlineRankApi)
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

// 总是会崩溃，不如直接重启
func supervisor(busChan chan DanmuMsg, roomInfoChan chan RoomInfo) {
	busChan <- DanmuMsg{
		Author:  "system",
		Content: "弹幕服务器初始化中",
		Type:    "NOTICE_MSG",
		Time:    time.Now(),
	}

	dc := DanmuClient{
		roomID:        uint32(config.Config.RoomId),
		auth:          config.Auth,
		conn:          new(websocket.Conn),
		unzlibChannel: make(chan []byte, 100),
	}

	defer func() {
		busChan <- DanmuMsg{
			Author:  "system",
			Content: "弹幕服务器已断开，正在重连",
			Type:    "NOTICE_MSG",
			Time:    time.Now(),
		}
		dc.isClosed = true
		dc.conn.Close()
		time.Sleep(1 * time.Second)
		supervisor(busChan, roomInfoChan)
	}()

	err := dc.connect()
	if err != nil {
		panic(err)
	}

	go dc.getHistory(busChan)
	go dc.receiveRawMsg(busChan)
	go dc.syncRoomInfo(roomInfoChan)
	go dc.heartBeat(busChan)

	for {
		time.Sleep(1 * time.Second)
		if dc.isClosed {
			return
		}
	}
}

func Run(busChan chan DanmuMsg, roomInfoChan chan RoomInfo) {
	go supervisor(busChan, roomInfoChan)
}
