package main

import (
	"bili/getter"
	"bili/sender"
	"bili/ui"
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Cookie string
	RoomId int64
}

var config Config

func init() {
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Printf("Error decoding config.toml: %s\n", err)
	}
}

func main() {
	busChan := make(chan string, 100)
	getter.Run(config.RoomId, busChan)
	sender.Run(config.Cookie)
	ui.Run(config.RoomId, busChan)
}
