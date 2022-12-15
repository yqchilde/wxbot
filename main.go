package main

import (
	"github.com/spf13/viper"
	"github.com/yqchilde/pkgs/log"
	"github.com/yqchilde/wxbot/engine/robot"
	"github.com/yqchilde/wxbot/framework/qianxun"
	"github.com/yqchilde/wxbot/framework/vlw"

	// 导入插件
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

func main() {
	// 初始化配置
	v := viper.New()
	v.SetConfigFile("config.yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("[main] 读取配置文件失败: %s", err.Error())
	}
	var c robot.Config
	if err := v.Unmarshal(&c); err != nil {
		log.Fatalf("[main] 解析配置文件失败: %s", err.Error())
	}

	// 初始化机器人
	if v.GetString("frameworks.name") == "" {
		log.Fatalf("[main] 未配置机器人框架")
	}
	switch v.GetString("frameworks.name") {
	case "qianxun":
		c.Framework = robot.IFramework(qianxun.New(
			v.GetString("botWxId"),
			v.GetString("frameworks.apiUrl"),
			v.GetString("frameworks.apiToken"),
			v.GetUint("frameworks.servePort"),
		))
	case "vlw":
		c.Framework = robot.IFramework(vlw.New(
			v.GetString("botWxId"),
			v.GetString("frameworks.apiUrl"),
			v.GetString("frameworks.apiToken"),
			v.GetUint("frameworks.servePort"),
		))
	default:
		log.Fatalf("[main] 未知机器人框架: %s", v.GetString("frameworks.name"))
	}
	robot.Run(&c)
}
