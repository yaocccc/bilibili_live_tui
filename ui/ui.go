package ui

import (
	"bili/config"
	"bili/getter"
	"bili/ui/theme1"
	"bili/ui/theme2"
	"bili/ui/theme3"
	"bili/ui/theme4"
)

func Run(busChan chan getter.DanmuMsg, roomInfoChan chan getter.RoomInfo) {
	switch config.Config.Theme {
	case 1: // theme1
		theme1.Run(busChan, roomInfoChan) // chat room
	case 2: // theme2
		theme2.Run(busChan, roomInfoChan) // pure
	case 3: // theme3
		theme3.Run(busChan, roomInfoChan) // simple
	case 4:
		theme4.Run(busChan, roomInfoChan) // info
	default:
		theme1.Run(busChan, roomInfoChan) // default theme1
	}
}
