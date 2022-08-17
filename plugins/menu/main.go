package menu

import (
	"strings"

	"github.com/eatmoreapple/openwechat"

	"github.com/yqchilde/wxbot/engine"
)

type Menu struct{}

var _ = engine.InstallPlugin(&Menu{})

func (m *Menu) OnRegister(event any) {}

func (m *Menu) OnEvent(event any) {
	if event != nil {
		msg := event.(*openwechat.Message)
		if msg.IsText() && msg.Content == "/menu" {
			if msg.IsSendByFriend() {
				reply := `YY Bot🤖
				🚀 输入 /img => 10s内发送表情获取表情原图
				🚀 输入 /plmm => 获取漂亮妹妹
				🚀 输入 /myb => 获取摸鱼办日记
				🚀 输入 /?? 拼音缩写 => 获取拼音缩写翻译
				🚀 输入 /kfc => 获取肯德基疯狂星期四骚话`
				msg.ReplyText(strings.ReplaceAll(reply, "\t", ""))
			} else if msg.IsSendByGroup() {
				reply := `YY Bot🤖
				🚀 输入 /img => 10s内发送表情获取表情原图
				🚀 输入 /plmm => 获取漂亮妹妹
				🚀 输入 /myb => 获取摸鱼办日记
				🚀 输入 /?? 拼音缩写 => 获取拼音缩写翻译
				🚀 输入 /kfc => 获取肯德基疯狂星期四骚话`
				msg.ReplyText(strings.ReplaceAll(reply, "\t", ""))
			}
		}
	}
}
