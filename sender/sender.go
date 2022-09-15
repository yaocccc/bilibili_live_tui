package sender

import (
	"fmt"
	"os"
	"strings"
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
	}
	time.AfterFunc(time.Second*10, heartbeat)
}

func Run(cookie string) {
	attrs := strings.Split(cookie, ";")
	kvs := make(map[string]string)
	for _, attr := range attrs {
		kv := strings.Split(attr, "=")
		k := strings.Trim(kv[0], " ")
		v := strings.Trim(kv[1], " ")
		kvs[k] = v
	}
	var auth bg.CookieAuth
	auth.SESSDATA = kvs["SESSDATA"]
	auth.DedeUserID = kvs["DedeUserID"]
	auth.DedeUserIDCkMd5 = kvs["DedeUserID__ckMd5"]
	auth.BiliJCT = kvs["bili_jct"]

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

func SendMsg(roomId int64, msg string) {
	r := []string{}
	for i := 0; i < len(msg); i += 20 {
		if i+20 < len(msg) {
			r = append(r, msg[i:i+20])
		} else {
			r = append(r, msg[i:])
		}
	}
	for _, v := range r {
		b.LiveSendDanmaku(roomId, 16777215, 25, 1, v, 0)
	}
}
