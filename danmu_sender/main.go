package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	bg "github.com/iyear/biligo"
)

type Config struct {
	ServerPort int
	Cookie     string
	RoomId     int64
}

var config Config
var auth bg.CookieAuth
var b *bg.BiliClient
var err error

func heartbeat() {
	start := time.Now()
	err := b.VideoHeartBeat(242531611, 173439442, int64(time.Since(start).Seconds()))
	if err != nil {
		log.Println("failed to send heartbeat; error:", err)
	}
	time.AfterFunc(time.Second*10, heartbeat)
}

// 初始化配置
func initConfig() {
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Printf("failed to parse config file; error: %v", err)
		os.Exit(0)
	}
}

// 初始化哔哩哔哩客户端
func initBilibili() {
	attrs := strings.Split(config.Cookie, ";")
	kvs := make(map[string]string)
	for _, attr := range attrs {
		kv := strings.Split(attr, "=")
		k := strings.Trim(kv[0], " ")
		v := strings.Trim(kv[1], " ")
		kvs[k] = v
	}
	auth.SESSDATA = kvs["SESSDATA"]
	auth.DedeUserID = kvs["DedeUserID"]
	auth.DedeUserIDCkMd5 = kvs["DedeUserID__ckMd5"]
	auth.BiliJCT = kvs["bili_jct"]

	b, err = bg.NewBiliClient(&bg.BiliSetting{
		Auth:      &auth,
		DebugMode: true,
	})
	if err != nil {
		fmt.Printf("failed to make new bili client; error: %v", err)
		os.Exit(0)
	}
	go heartbeat()
}

func reqHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	msg := string(bodyBytes)
	err = b.LiveSendDanmaku(config.RoomId, 16777215, 25, 1, msg, 0)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	fmt.Fprintln(w, "ok")
}

func init() {
	initConfig()
	initBilibili()
}

func main() {
	http.HandleFunc("/", reqHandler)
	http.ListenAndServe("localhost:9527", nil)
}
