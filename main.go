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

var ping = true

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
			if ping = net.PingConn(ipPort, time.Second*20); !ping {
				log.Warn("[main] 无法连接到千寻框架，网络无法Ping通")
			}
		}
	case "VLW", "vlw":
		f = robot.IFramework(vlw.New(c.ServerPort, c.BotWxId, c.Framework.ApiUrl, c.Framework.ApiToken))
		if ipPort, err := net.CheckoutIpPort(c.Framework.ApiUrl); err == nil {
			if ping = net.PingConn(ipPort, time.Second*20); !ping {
				log.Warn("[main] 无法连接到VLW框架，网络无法Ping通")
			}
		}
	default:
		log.Fatalf("[main] 请在配置文件中指定机器人框架后再启动")
	}

	robot.WxBot = robot.Init(c, f)
	if !ping {
		log.Println("[main] 开始获取账号数据...")
		friendsList, err := robot.WxBot.Framework.GetFriendsList(true)
		if err != nil {
			log.Errorf("[main] 获取好友列表失败，error: %s", err.Error())
		}
		groupList, err := robot.WxBot.Framework.GetGroupList(true)
		if err != nil {
			log.Errorf("[main] 获取群组列表失败，error: %s", err.Error())
		}
		subscriptionList, err := robot.WxBot.Framework.GetSubscriptionList(true)
		if err != nil {
			log.Errorf("[main] 获取公众号列表失败，error: %s", err.Error())
		}
		robot.WxBot.FriendsList = friendsList
		robot.WxBot.GroupList = groupList
		robot.WxBot.SubscriptionList = subscriptionList
		log.Printf("[main] 共获取到%d个好友", len(friendsList))
		log.Printf("[main] 共获取到%d个群组", len(groupList))
		log.Printf("[main] 共获取到%d个公众号", len(subscriptionList))
	}

	log.Printf("[main] 机器人%s开始工作", c.BotNickname)
	robot.WxBot.Run()
}
