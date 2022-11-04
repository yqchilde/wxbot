package zaobao

import (
	"github.com/yqchilde/wxbot/engine"
	"github.com/yqchilde/wxbot/engine/robot"
)

type ZaoBao struct{ engine.PluginMagic }

var (
	pluginInfo = &ZaoBao{
		engine.PluginMagic{
			Desc:     "🚀 输入 {每日早报|早报} => 获取每天60s读懂世界",
			Commands: []string{"每日早报", "早报"},
		},
	}
	_ = engine.InstallPlugin(pluginInfo)
)

func (p *ZaoBao) OnRegister() {}

func (p *ZaoBao) OnEvent(msg *robot.Message) {
	if msg.MatchTextCommand(pluginInfo.Commands) {
		msg.ReplyImage("https://api.qqsuu.cn/api/dm-60s?type=image")
	}
}
