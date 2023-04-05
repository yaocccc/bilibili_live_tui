package main

import (
	"bili/config"
	"bili/getter"
	"bili/sender"
	"bili/ui"
	"os"
	"os/exec"
	"strings"
)

/** 用于修正环境变量 */
func fixCharset() {
	locale := os.Getenv("LANG")
	var asianCharset bool
	var wideCharset = []string{"zh_", "jp_", "ko_", "ja_", "th_", "hi_"}
	for k := range wideCharset {
		if strings.HasPrefix(locale, wideCharset[k]) {
			asianCharset = true
		}
	}
	if asianCharset {
		os.Setenv("LANG", "C.UTF-8")
		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		os.Exit(0)
	}
}

func main() {
	fixCharset()
	config.Init()
	busChan := make(chan getter.DanmuMsg, 100)
	roomInfoChan := make(chan getter.RoomInfo, 100)
	getter.Run(busChan, roomInfoChan)
	sender.Run()
	ui.Run(busChan, roomInfoChan)
}
