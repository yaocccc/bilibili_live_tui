package ui

import (
	"bili/config"
	"bili/getter"
	"bili/ui/theme1"
	"bili/ui/theme2"
	"bili/ui/theme3"
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

func Run(busChan chan getter.DanmuMsg, roomInfoChan chan getter.RoomInfo) {
	fixCharset()
	switch config.Config.Theme {
	case 1: // theme1
		theme1.Run(busChan, roomInfoChan) // chat room
	case 2: // theme2
		theme2.Run(busChan, roomInfoChan) // pure
	case 3: // theme3
		theme3.Run(busChan, roomInfoChan) // simple
	default:
		theme1.Run(busChan, roomInfoChan) // default theme1
	}
}
