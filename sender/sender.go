package sender

import (
	"fmt"
	"os"
	"time"

	bg "github.com/iyear/biligo"
)

var b *bg.BiliClient
var err error

func heartbeat() {
	start := time.Now()
	err := b.VideoHeartBeat(242531611, 173439442, int64(time.Since(start).Seconds()))
	if err != nil {
		fmt.Println("failed to send heartbeat; error:", err)
		os.Exit(0)
	}
	time.AfterFunc(time.Second*10, heartbeat)
}

func SendMsg(roomId int64, msg string, busChan chan []string) {
	msgRune := []rune(msg)
	for i := 0; i < len(msgRune); i += 20 {
		err = nil
		if i+20 < len(msgRune) {
			err = b.LiveSendDanmaku(roomId, 16777215, 25, 1, string(msgRune[i:i+20]), 0)
			time.Sleep(time.Second * 1)
		} else {
			err = b.LiveSendDanmaku(roomId, 16777215, 25, 1, string(msgRune[i:]), 0)
		}
		if err != nil {
			busChan <- []string{"error", err.Error()}
		}
	}
}

func Run(auth bg.CookieAuth) {
	b, err = bg.NewBiliClient(&bg.BiliSetting{
		Auth:      &auth,
		DebugMode: false,
	})
	if err != nil {
		fmt.Printf("failed to make new bili client; error: %v", err)
		os.Exit(0)
	}
	go heartbeat()
}
