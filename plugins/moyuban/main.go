package moyuban

import (
	"github.com/imroc/req/v3"

	"github.com/yqchilde/wxbot/engine"
	"github.com/yqchilde/wxbot/engine/robot"
)

type MoYuBan struct{ engine.PluginMagic }

var (
	pluginInfo = &MoYuBan{
		engine.PluginMagic{
			Desc:     "🚀 输入 {摸鱼日记} => 获取摸鱼办日记",
			Commands: []string{"摸鱼日记", "摸鱼"},
		},
	}
	_ = engine.InstallPlugin(pluginInfo)
)

func (m *MoYuBan) OnRegister() {}

func (m *MoYuBan) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.MatchTextCommand(pluginInfo.Commands) {
			if url := getMoYuData(); url != "" {
				msg.ReplyImage(url)
			} else {
				msg.ReplyText("获取摸鱼办日记失败")
			}
		}
	}
}

func getMoYuData() (url string) {
	type Resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			MoyuUrl string `json:"moyu_url"`
		} `json:"data"`
	}
	var resp Resp
	if err := req.C().Get("https://api.j4u.ink/v1/store/other/proxy/remote/moyu.json").Do().Into(&resp); err != nil {
		return ""
	}
	if resp.Code != 200 || resp.Message != "success" {
		return ""
	}
	return resp.Data.MoyuUrl
}
