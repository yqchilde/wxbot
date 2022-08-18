package plmm

import (
	"os"

	"github.com/eatmoreapple/openwechat"

	"github.com/yqchilde/wxbot/engine"
)

type Plmm struct {
	engine.PluginMagic
	Enable    bool   `yaml:"enable"`
	Dir       string `yaml:"dir"`
	Url       string `yaml:"url"`
	AppId     string `yaml:"appId"`
	AppSecret string `yaml:"appSecret"`
}

var (
	pluginInfo = &Plmm{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 /plmm => 获取漂亮妹妹",
			Commands: []string{"/plmm"},
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (p *Plmm) OnRegister(event any) {
	err := os.MkdirAll(plugin.RawConfig.Get("dir").(string), os.ModePerm)
	if err != nil {
		panic("init plmm img dir error: " + err.Error())
	}
}

func (p *Plmm) OnEvent(event any) {
	if event != nil {
		msg := event.(*openwechat.Message)
		if msg.IsText() && msg.Content == pluginInfo.Commands[0] {
			getPlmmPhoto(msg)
		}
	}
}
