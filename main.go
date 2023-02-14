package main

import (
	"time"

	"github.com/spf13/viper"
	"github.com/yqchilde/pkgs/net"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/robot"
	"github.com/yqchilde/wxbot/framework/qianxun"
	"github.com/yqchilde/wxbot/framework/vlw"

	// 导入插件, 变更插件请查看README
	_ "github.com/yqchilde/wxbot/engine/plugins"
)

func main() {
	v := viper.New()
	v.SetConfigFile("config.yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("[main] 读取配置文件失败: %s", err.Error())
	}
	c := new(robot.Config)
	if err := v.Unmarshal(c); err != nil {
		log.Fatalf("[main] 解析配置文件失败: %s", err.Error())
	}

	f := robot.IFramework(nil)
	switch c.Framework.Name {
	case "千寻", "qianxun":
		f = robot.IFramework(qianxun.New(c.ServerPort, c.BotWxId, c.Framework.ApiUrl, c.Framework.ApiToken))
		if ipPort, err := net.CheckoutIpPort(c.Framework.ApiUrl); err == nil {
			if ping := net.PingConn(ipPort, time.Second*20); !ping {
				log.Warn("[main] 无法连接到千寻框架，网络无法Ping通")
			}
		}
	case "VLW", "vlw":
		f = robot.IFramework(vlw.New(c.ServerPort, c.BotWxId, c.Framework.ApiUrl, c.Framework.ApiToken))
		if ipPort, err := net.CheckoutIpPort(c.Framework.ApiUrl); err == nil {
			if ping := net.PingConn(ipPort, time.Second*20); !ping {
				log.Warn("[main] 无法连接到VLW框架，网络无法Ping通")
			}
		}
	default:
		log.Fatalf("[main] 请在配置文件中指定机器人框架后再启动")
	}

	robot.Run(c, f)
}
