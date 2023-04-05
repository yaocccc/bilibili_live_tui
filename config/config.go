package config

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/BurntSushi/toml"
	bg "github.com/iyear/biligo"
)

type ConfigType struct {
	Cookie       string // 登录cookie
	RoomId       int64  // 直播间id
	Theme        int64  // 主题
	SingleLine   int64  // 是否开启单行
	ShowTime     int64  // 是否显示时间
	TimeColor    string // 时间颜色
	NameColor    string // 名字颜色
	ContentColor string // 内容颜色
	FrameColor   string // 边框颜色
	InfoColor    string // 房间信息颜色
	RankColor    string // 排行榜颜色
}

var Auth bg.CookieAuth
var Config ConfigType

func Init() {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	homeDir := currentUser.HomeDir
	path := homeDir + "/.config/bili"
	if err = os.MkdirAll(path, 0755); err != nil {
		panic(err)
	}
	configFile := path + "/config.toml"
	_, err = os.Stat(configFile)
	if os.IsNotExist(err) {
		config := ConfigType{
			Cookie:       "从你BILIBILI的请求里抓一个Cookie",
			RoomId:       23333333,
			Theme:        1,
			SingleLine:   1,
			ShowTime:     1,
			TimeColor:    "#FFFFFF",
			NameColor:    "#FFFFFF",
			ContentColor: "#FFFFFF",
			FrameColor:   "#FFFFFF",
			InfoColor:    "#FFFFFF",
			RankColor:    "#FFFFFF",
		}
		f, err := os.Create(configFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		if err := toml.NewEncoder(f).Encode(config); err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	if _, err := toml.DecodeFile(configFile, &Config); err != nil {
		fmt.Printf("Error decoding config.toml: %s\n", err)
	}
	roomId := int64(-1)
	theme := int64(-1)
	single_line := int64(-1)
	show_time := int64(-1)
	flag.StringVar(&configFile, "c", "config.toml", "usage for config")
	flag.Int64Var(&roomId, "r", -1, "usage for room id")
	flag.Int64Var(&theme, "t", -1, "usage for theme")
	flag.Int64Var(&single_line, "l", -1, "usage for single_line")
	flag.Int64Var(&show_time, "s", -1, "usage for show_time")
	flag.Parse()

	if roomId != -1 {
		Config.RoomId = roomId
	}
	if theme != -1 {
		Config.Theme = theme
	}
	if single_line != -1 {
		Config.SingleLine = single_line
	}
	if show_time != -1 {
		Config.ShowTime = show_time
	}
	if Config.TimeColor == "" {
		Config.TimeColor = "#bbbbbb"
	}
	if Config.NameColor == "" {
		Config.NameColor = "#bbbbbb"
	}
	if Config.ContentColor == "" {
		Config.ContentColor = "#bbbbbb"
	}
	if Config.TimeColor == "" {
		Config.TimeColor = "#bbbbbb"
	}
	if Config.NameColor == "" {
		Config.NameColor = "#bbbbbb"
	}
	if Config.ContentColor == "" {
		Config.ContentColor = "#bbbbbb"
	}
	if Config.InfoColor == "" {
		Config.InfoColor = "#bbbbbb"
	}
	if Config.RankColor == "" {
		Config.RankColor = "#bbbbbb"
	}
	if Config.FrameColor == "" {
		Config.FrameColor = "#bbbbbb"
	}

	attrs := strings.Split(Config.Cookie, ";")
	kvs := make(map[string]string)
	for _, attr := range attrs {
		kv := strings.Split(attr, "=")
		k := strings.Trim(kv[0], " ")
		v := strings.Trim(kv[1], " ")
		kvs[k] = v
	}
	Auth.SESSDATA = kvs["SESSDATA"]
	Auth.DedeUserID = kvs["DedeUserID"]
	Auth.DedeUserIDCkMd5 = kvs["DedeUserID__ckMd5"]
	Auth.BiliJCT = kvs["bili_jct"]
}
