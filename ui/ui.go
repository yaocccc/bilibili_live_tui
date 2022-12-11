package ui

import (
	"bili/getter"
	"bili/ui/theme1"
	"bili/ui/theme2"
)

func Run(roomId int64, theme int64, busChan chan getter.DanmuMsg, roomInfoChan chan getter.RoomInfo) {
	switch theme {
	case 1: // theme1
		theme1.Run(roomId, busChan, roomInfoChan) // chat room
	case 2: // theme2
		theme2.Run(roomId, busChan, roomInfoChan) // simple
	default:
		theme1.Run(roomId, busChan, roomInfoChan) // default theme1
	}
}
