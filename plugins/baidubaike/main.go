package baidubaike

import (
	"fmt"

	"github.com/imroc/req/v3"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/robot"
)

func init() {
	engine := control.Register("baidubaike", &control.Options{
		Alias: "百度百科",
		Help: "指令:\n" +
			"* 百度百科 [查询内容]",
	})

	engine.OnRegex(`^百度百科 ?(.*?)$`).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		word := ctx.State["regex_matched"].([]string)[1]
		if data, err := getBaiKe(word); err == nil {
			if data == nil {
				ctx.ReplyText("没查到该百科含义")
			} else {
				ctx.ReplyText("🏷" + word + ":\n" + fmt.Sprintf("%s\n🔎 摘要: %s\n©️ 版权: %s", data.Desc, data.Abstract, data.Copyrights))
			}
		} else {
			ctx.ReplyText("查询失败，这一定不是bug🤔")
		}
	})
}

type apiResponse struct {
	Key        string `json:"key"`
	Desc       string `json:"desc"`
	Abstract   string `json:"abstract"`
	Copyrights string `json:"copyrights"`
}

func getBaiKe(keyword string) (*apiResponse, error) {
	var data apiResponse
	api := "https://baike.baidu.com/api/openapi/BaikeLemmaCardApi?appid=379020&bk_length=1000&bk_key=" + keyword
	if err := req.C().Get(api).Do().Into(&data); err != nil {
		return nil, err
	}
	if len(data.Abstract) == 0 {
		return nil, nil
	}
	return &data, nil
}
