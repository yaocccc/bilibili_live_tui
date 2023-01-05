package config

import (
	"flag"
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	bg "github.com/iyear/biligo"
)

type ConfigType struct {
	Cookie       string
	RoomId       int64
	Theme        int64
	TimeColor    string
	NameColor    string
	ContentColor string
}

var Config ConfigType
var Auth bg.CookieAuth

func Init() {
	configFile := ""
	roomId := int64(0)
	theme := int64(0)
	flag.StringVar(&configFile, "c", "config.toml", "usage for config")
	flag.Int64Var(&roomId, "r", 0, "usage for room id")
	flag.Int64Var(&theme, "t", 0, "usage for theme")
	flag.Parse()

	if _, err := toml.DecodeFile(configFile, &Config); err != nil {
		fmt.Printf("Error decoding config.toml: %s\n", err)
	}

	if roomId != 0 {
		Config.RoomId = roomId
	}
	if theme != 0 {
		Config.Theme = theme
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
