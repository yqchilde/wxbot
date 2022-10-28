package jingdong

import (
	"github.com/yqchilde/wxbot/engine"
	"github.com/yqchilde/wxbot/engine/robot"
)

type JingDong struct{ engine.PluginMagic }

var (
	pluginInfo = &JingDong{
		engine.PluginMagic{
			Desc:     "🚀 输入 {京东上车} => 快上车和我一起挂京豆",
			Commands: []string{"京东上车"},
		},
	}
	_ = engine.InstallPlugin(pluginInfo)
)

func (p *JingDong) OnRegister() {}

func (p *JingDong) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.MatchTextCommand(pluginInfo.Commands) {
			if err := msg.ReplyImage("C:\\Users\\Administrator\\Pictures\\jd\\qrcode.png"); err != nil {
				msg.ReplyText(err.Error())
			}
		}
	}
}
