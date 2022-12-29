package main

import (
	"time"

	"github.com/spf13/viper"
	"github.com/yqchilde/pkgs/net"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/robot"
	"github.com/yqchilde/wxbot/framework/qianxun"
	"github.com/yqchilde/wxbot/framework/vlw"

	// 导入插件, 不需要的插件可以注释掉或者删除
	_ "github.com/yqchilde/wxbot/plugins/baidubaike"   // 百度百科
	_ "github.com/yqchilde/wxbot/plugins/chatgpt"      // GPT聊天
	_ "github.com/yqchilde/wxbot/plugins/crazykfc"     // 肯德基疯狂星期四骚话
	_ "github.com/yqchilde/wxbot/plugins/jingdong"     // 京豆上车
	_ "github.com/yqchilde/wxbot/plugins/manager"      // 群组管理相关
	_ "github.com/yqchilde/wxbot/plugins/memepicture"  // 表情包原图
	_ "github.com/yqchilde/wxbot/plugins/moyuban"      // 摸鱼办
	_ "github.com/yqchilde/wxbot/plugins/pinyinsuoxie" // 拼音缩写翻译
	_ "github.com/yqchilde/wxbot/plugins/plmm"         // 漂亮妹妹
	_ "github.com/yqchilde/wxbot/plugins/weather"      // 天气查询
	_ "github.com/yqchilde/wxbot/plugins/zaobao"       // 每日早报
)

var conf robot.Config

func main() {
	// 初始化配置
	v := viper.New()
	v.SetConfigFile("config.yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("[main] 读取配置文件失败: %s", err.Error())
	}
	if err := v.Unmarshal(&conf); err != nil {
		log.Fatalf("[main] 解析配置文件失败: %s", err.Error())
	}

	// 初始化机器人
	switch v.GetString("frameworks.name") {
	case "千寻", "qianxun":
		conf.Framework = robot.IFramework(qianxun.New(
			v.GetString("botWxId"),
			v.GetString("frameworks.apiUrl"),
			v.GetString("frameworks.apiToken"),
			v.GetUint("frameworks.servePort"),
		))
		if ipPort, err := net.CheckoutIpPort(v.GetString("frameworks.apiUrl")); err == nil {
			if ping := net.PingConn(ipPort, time.Second*20); !ping {
				log.Fatalf("[main] 无法连接到千寻框架，网络无法Ping通")
			}
		}
	case "VLW", "vlw":
		conf.Framework = robot.IFramework(vlw.New(
			v.GetString("botWxId"),
			v.GetString("frameworks.apiUrl"),
			v.GetString("frameworks.apiToken"),
			v.GetUint("frameworks.servePort"),
		))
		if ipPort, err := net.CheckoutIpPort(v.GetString("frameworks.apiUrl")); err == nil {
			if ping := net.PingConn(ipPort, time.Second*20); !ping {
				log.Fatalf("[main] 无法连接到VLW框架，网络无法Ping通")
			}
		}
	default:
		log.Fatalf("[main] 请在配置文件中指定机器人框架后再启动")
	}
	robot.Run(&conf)
}
