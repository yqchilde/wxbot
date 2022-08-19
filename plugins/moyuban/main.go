package moyuban

import (
	"embed"

	"github.com/yqchilde/wxbot/engine"
	"github.com/yqchilde/wxbot/engine/robot"
)

type MoYuBan struct{ engine.PluginMagic }

var (
	pluginInfo = &MoYuBan{
		engine.PluginMagic{
			Desc:     "🚀 输入 /myb => 获取摸鱼办日记",
			Commands: []string{"/myb"},
		},
	}
	_ = engine.InstallPlugin(pluginInfo)
)

//go:embed holiday.json
var f embed.FS

func (m *MoYuBan) OnRegister() {}

func (m *MoYuBan) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.IsText() && msg.Content == pluginInfo.Commands[0] {
			if notes, err := DailyLifeNotes(""); err == nil {
				msg.ReplyText(notes)
			} else {
				msg.ReplyText("查询失败，这一定不是bug🤔")
			}
		}
	}
}
