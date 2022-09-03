package plmm

import (
	"os"

	"github.com/yqchilde/wxbot/engine"
	"github.com/yqchilde/wxbot/engine/robot"
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
			Desc:     "🚀 输入 {漂亮妹妹} => 获取漂亮妹妹",
			Commands: []string{"漂亮妹妹"},
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (p *Plmm) OnRegister() {
	err := os.MkdirAll(plugin.RawConfig.Get("dir").(string), os.ModePerm)
	if err != nil {
		panic("init plmm img dir error: " + err.Error())
	}
}

func (p *Plmm) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.MatchTextCommand(plugin.Commands) {
			getPlmmPhoto(msg)
		}
	}
}
