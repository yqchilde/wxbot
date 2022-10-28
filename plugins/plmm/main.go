package plmm

import (
	"fmt"
	"os"

	"github.com/imroc/req/v3"
	"github.com/yqchilde/pkgs/log"

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

var plmmUrlStorage []string

func getPlmmPhoto(msg *robot.Message) {
	var plmmConf Plmm
	plugin.RawConfig.Unmarshal(&plmmConf)

	if len(plmmUrlStorage) > 50 {
		if err := msg.ReplyImage(plmmUrlStorage[0]); err != nil {
			msg.ReplyText(err.Error())
		}
		plmmUrlStorage = plmmUrlStorage[1:]
	} else {
		var resp PlmmApiResponse
		api := fmt.Sprintf("%s?app_id=%s&app_secret=%s", plmmConf.Url, plmmConf.AppId, plmmConf.AppSecret)
		err := req.C().SetBaseURL(api).Get().Do().Into(&resp)
		if err != nil {
			log.Errorf("get plmm photo error: %s", err.Error())
			return
		}
		if resp.Code != 1 {
			log.Errorf("getPlmmPhoto api error: %v", resp.Msg)
			return
		}
		for _, val := range resp.Data {
			plmmUrlStorage = append(plmmUrlStorage, val.ImageUrl)
		}
		if err := msg.ReplyImage(plmmUrlStorage[0]); err != nil {
			msg.ReplyText(err.Error())
		}
		plmmUrlStorage = plmmUrlStorage[1:]
	}
}
