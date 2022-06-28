package main

import (
	"danmu_geter/danmu"
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	RoomId uint32
}

var config Config
var dc *danmu.DanmuClient

func init() {
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Printf("Error decoding config.toml: %s\n", err)
	}

	dc = danmu.NewDanmuClient(config.RoomId)
}

func main() {
	dc.Run()
	time.Sleep(time.Hour * 999)
}
