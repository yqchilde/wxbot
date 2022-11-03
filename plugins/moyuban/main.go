package moyuban

import (
	"github.com/yqchilde/wxbot/engine"
	"github.com/yqchilde/wxbot/engine/robot"
)

type MoYuBan struct{ engine.PluginMagic }

var (
	pluginInfo = &MoYuBan{
		engine.PluginMagic{
			Desc:     "🚀 输入 {摸鱼日历|摸鱼} => 获取摸鱼办日历",
			Commands: []string{"摸鱼日历", "摸鱼"},
			Weight:   97,
		},
	}
	_ = engine.InstallPlugin(pluginInfo)
)

func (m *MoYuBan) OnRegister() {}

func (m *MoYuBan) OnEvent(msg *robot.Message) {
	if msg.MatchTextCommand(pluginInfo.Commands) {
		msg.ReplyImage("https://api.vvhan.com/api/moyu")
	}
}
