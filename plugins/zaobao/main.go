package zaobao

import (
	"github.com/imroc/req/v3"

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
		if zaoBao, err := getZaoBao(); err == nil {
			msg.ReplyImage(zaoBao)
		}
	}
}

func getZaoBao() (string, error) {
	type Resp struct {
		Code     int    `json:"code"`
		Msg      string `json:"msg"`
		ImageUrl string `json:"imageUrl"`
		Datatime string `json:"datatime"`
	}
	var resp Resp
	if err := req.C().Get("http://dwz.2xb.cn/zaob").Do().Into(&resp); err != nil {
		return "", err
	}
	return resp.ImageUrl, nil
}
