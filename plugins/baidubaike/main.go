package baidubaike

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/yqchilde/wxbot/engine"
	"github.com/yqchilde/wxbot/engine/robot"
)

type BaiDuBaiKe struct{ engine.PluginMagic }

var (
	pluginInfo = &BaiDuBaiKe{
		engine.PluginMagic{
			Desc:     "🚀 输入 {百度百科 XX} => 获取百度百科解释，Ps:百度百科 okr",
			Commands: []string{"^百度百科 ?(.*?)$"},
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (p *BaiDuBaiKe) OnRegister() {}

func (p *BaiDuBaiKe) OnEvent(msg *robot.Message) {
	if idx, ok := msg.MatchRegexCommand(pluginInfo.Commands); ok {
		var re = regexp.MustCompile(pluginInfo.Commands[idx])
		match := re.FindAllStringSubmatch(msg.Content.Msg, -1)
		if len(match) > 0 && len(match[0]) > 1 {
			if data, err := getBaiKe(match[0][1]); err == nil {
				if data == nil {
					msg.ReplyText("没查到该百科含义")
				} else {
					msg.ReplyText("🏷️ " + match[0][1] + ": " + fmt.Sprintf("%s\n🔎 摘要: %s\n© 版权: %s", data.Desc, data.Abstract, data.Copyrights))
				}
			} else {
				msg.ReplyText("查询失败，这一定不是bug🤔")
			}
		}
	}
}

func getBaiKe(keyword string) (*ApiResponse, error) {
	api := "https://baike.baidu.com/api/openapi/BaikeLemmaCardApi?appid=379020&bk_length=600&bk_key=" + keyword
	resp, err := http.Get(api)
	if err != nil {
		plugin.Errorf("failed to get baike api, err: %v", err)
		return nil, err
	}
	readAll, err := io.ReadAll(resp.Body)
	if err != nil {
		plugin.Errorf("failed to read resp body, err: %v", err)
		return nil, err
	}
	var data ApiResponse
	if err := json.Unmarshal(readAll, &data); err != nil {
		plugin.Errorf("failed to unmarshal api response, err: %v", err)
		return nil, err
	}
	if len(data.Abstract) == 0 {
		return nil, nil
	}
	return &data, nil
}
