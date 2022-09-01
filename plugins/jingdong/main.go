package jingdong

import (
	"os"
	"strings"

	"github.com/yqchilde/pkgs/log"

	"github.com/yqchilde/wxbot/engine"
	"github.com/yqchilde/wxbot/engine/robot"
)

type JingDong struct{ engine.PluginMagic }

var (
	pluginInfo = &JingDong{
		engine.PluginMagic{
			Desc:     "🚀 输入 {/jd} => 获取京东提交CK二维码",
			Commands: []string{"/jd"},
		},
	}
	_ = engine.InstallPlugin(pluginInfo)
)

func (p *JingDong) OnRegister() {}

func (p *JingDong) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.MatchTextCommand(pluginInfo.Commands) {
			img, err := os.Open("./imgs/jingdong/qrcode.png")
			if err != nil {
				msg.ReplyText("Err: " + err.Error())
			}
			defer img.Close()

			if _, err := msg.ReplyImage(img); err != nil {
				if strings.Contains(err.Error(), "operate too often") {
					msg.ReplyText("Warn: 被微信ban了，请稍后再试")
				} else {
					log.Errorf("msg.ReplyImage reply image error: %v", err)
				}
			}
		}
	}
}
