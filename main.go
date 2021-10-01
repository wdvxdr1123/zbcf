package main

import (
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"

	_ "github.com/wdvxdr1123/zbcf/plugins/contest"
	_ "github.com/wdvxdr1123/zbcf/plugins/user"
	"github.com/wdvxdr1123/zbcf/util"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(util.LogFormat{})
	zero.Run(zero.Config{
		NickName:      []string{"bot"},
		CommandPrefix: "/",
		SuperUsers:    []string{"123456"},
		Driver: []zero.Driver{
			driver.NewWebSocketClient("ws://127.0.0.1:6700/", ""),
		},
	})
	select {}
}
