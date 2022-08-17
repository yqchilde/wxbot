package moyuban

import (
	"embed"

	"github.com/eatmoreapple/openwechat"

	"github.com/yqchilde/wxbot/engine"
)

type MoYuBan struct{}

var _ = engine.InstallPlugin(&MoYuBan{})

//go:embed holiday.json
var f embed.FS

func (m *MoYuBan) OnRegister(event any) {}

func (m *MoYuBan) OnEvent(event any) {
	if event != nil {
		msg := event.(*openwechat.Message)
		if msg.IsText() && msg.Content == "/myb" {
			if notes, err := DailyLifeNotes(""); err == nil {
				msg.ReplyText(notes)
			} else {
				msg.ReplyText("查询失败，这一定不是bug🤔")
			}
		}
	}
}
