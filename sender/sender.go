package sender

import (
	"bili/config"
	"bili/getter"
	"fmt"
	"os"
	"time"

	bg "github.com/iyear/biligo"
)

var bc *bg.BiliClient
var err error

func heartbeat() {
	start := time.Now()
	err := bc.VideoHeartBeat(242531611, 173439442, int64(time.Since(start).Seconds()))
	if err != nil {
		fmt.Println("failed to send heartbeat; error:", err)
		os.Exit(0)
	}
	time.AfterFunc(time.Second*10, heartbeat)
}

func SendMsg(roomId int64, msg string, busChan chan getter.DanmuMsg) {
	msgRune := []rune(msg)
	for i := 0; i < len(msgRune); i += 20 {
		err = nil
		if i+20 < len(msgRune) {
			err = bc.LiveSendDanmaku(roomId, 16777215, 25, 1, string(msgRune[i:i+20]), 0)
			time.Sleep(time.Second * 1)
		} else {
			err = bc.LiveSendDanmaku(roomId, 16777215, 25, 1, string(msgRune[i:]), 0)
		}
		if err != nil {
			busChan <- getter.DanmuMsg{Author: "system", Content: "发送弹幕失败", Type: ""}
		}
	}
}

func Run() {
	for retry := 0; retry < 3; retry++ {
		bc, err = bg.NewBiliClient(&bg.BiliSetting{
			Auth:      &config.Auth,
			DebugMode: false,
		})
		if err == nil {
			break
		}
		time.Sleep(time.Second * 1)
	}
	go heartbeat()
}
