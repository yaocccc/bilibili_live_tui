package main

import (
	"bili/config"
	"bili/getter"
	"bili/sender"
	"bili/ui"
)

func init() {
	config.Init()
}

func main() {
	busChan := make(chan getter.DanmuMsg, 100)
	roomInfoChan := make(chan getter.RoomInfo, 100)
	getter.Run(busChan, roomInfoChan)
	sender.Run()
	ui.Run(busChan, roomInfoChan)
}
