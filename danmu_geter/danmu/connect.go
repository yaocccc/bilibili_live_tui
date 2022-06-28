package danmu

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/asmcos/requests"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
)

var (
	getDanmuInfo = "https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo?id=%d&type=0"
)

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
func (d *DanmuClient) receiveRawMsg() {
	for {
		_, msg, _ := d.conn.ReadMessage()
		if msg[7] == 2 {
			msgs := splitMsg(zlibUnCompress(msg[16:]))
			for _, m := range msgs {
				d.unzlibChannel <- m
			}
		} else if msg[11] == 3 {
			d.heartBeatChannel <- msg
		} else {
			d.serverNoticeChannel <- msg
		}
	}
}

func (d *DanmuClient) Run() {
	d.connect()
	go d.process()
	go d.heartBeat()
	go d.receiveRawMsg()
}
