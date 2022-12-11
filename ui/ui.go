package ui

import (
	"bili/getter"
	"bili/ui/theme1"
	"bili/ui/theme2"
	"os"
	"os/exec"
	"strings"
)

var (
	wideCharset = []string{"zh_", "jp_", "ko_", "ja_", "th_", "hi_"}
)

func fixCharset() {
	locale := os.Getenv("LANG")

	var asianCharset bool
	for k := range wideCharset {
		if strings.HasPrefix(locale, wideCharset[k]) {
			asianCharset = true
		}
	}

	if asianCharset {
		os.Setenv("LANG", "C.UTF-8")
		cmd := exec.Command(os.Args[0])
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		os.Exit(0)
	}
}

func Run(roomId int64, theme int64, busChan chan getter.DanmuMsg, roomInfoChan chan getter.RoomInfo) {
	fixCharset()
	switch theme {
	case 1: // theme1
		theme1.Run(roomId, busChan, roomInfoChan) // chat room
	case 2: // theme2
		theme2.Run(roomId, busChan, roomInfoChan) // simple
	default:
		theme1.Run(roomId, busChan, roomInfoChan) // default theme1
	}
}
