package douyingirl

import (
	"github.com/yqchilde/wxbot/engine"
	"github.com/yqchilde/wxbot/engine/robot"
)

type DouYinGirl struct{ engine.PluginMagic }

var (
	pluginInfo = &DouYinGirl{
		engine.PluginMagic{
			Desc:     "🚀 输入 {抖音小姐姐} => 获取抖音小姐姐视频",
			Commands: []string{"抖音小姐姐"},
		},
	}
	_ = engine.InstallPlugin(pluginInfo)
)

func (p *DouYinGirl) OnRegister() {}

func (p *DouYinGirl) OnEvent(msg *robot.Message) {
	if msg.MatchTextCommand(pluginInfo.Commands) {
		msg.ReplyShareLink("抖音小姐姐", "每次点进来都不一样呦", "https://www.haofang365.com/uploads/20211114/zip_16368195136np7EF.jpg", "http://nvz.bcyle.com/999.php")
	}
}
