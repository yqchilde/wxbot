package plmm

import (
	"os"

	"github.com/eatmoreapple/openwechat"

	"github.com/yqchilde/wxbot/engine"
)

type Plmm struct {
	Enable    bool   `yaml:"enable"`
	Dir       string `yaml:"dir"`
	Url       string `yaml:"url"`
	AppId     string `yaml:"appId"`
	AppSecret string `yaml:"appSecret"`
}

var plugin = engine.InstallPlugin(&Plmm{})

func (p *Plmm) OnRegister(event any) {
	err := os.MkdirAll(plugin.RawConfig.Get("dir").(string), os.ModePerm)
	if err != nil {
		panic("init plmm img dir error: " + err.Error())
	}
}

func (p *Plmm) OnEvent(event any) {
	if event != nil {
		msg := event.(*openwechat.Message)
		if msg.IsText() && msg.Content == "/plmm" {
			getPlmmPhoto(msg)
		}
	}
}
