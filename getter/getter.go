package getter

import "github.com/gorilla/websocket"

type DanmuClient struct {
	roomID        uint32
	conn          *websocket.Conn
	unzlibChannel chan []byte
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

func Run(roomID int64, busChan chan []string) {
	dc := DanmuClient{
		roomID:        uint32(roomID),
		conn:          new(websocket.Conn),
		unzlibChannel: make(chan []byte, 100),
	}
	dc.Run(busChan)
}
