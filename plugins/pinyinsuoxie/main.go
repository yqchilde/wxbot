package pinyinsuoxie

import (
	"fmt"
	"strings"

	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/robot"
)

func init() {
	engine := control.Register("chasuoxie", &control.Options{
		Alias: "查缩写",
		Help: "描述:\n" +
			"奇奇怪怪的拼音缩写咱也不知道啥意思啊，快来查一查\n\n" +
			"指令:\n" +
			"* 查缩写 [内容] -> 获取拼音缩写翻译，Ps:查缩写 yyds",
	})
	engine.OnRegex(`^查缩写 ?([a-zA-Z0-9]+)$`).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		word := ctx.State["regex_matched"].([]string)[1]
		if data, err := transPinYinSuoXie(word); err == nil {
			if len(data) == 0 {
				ctx.ReplyText("没查到该缩写含义")
			} else {
				ctx.ReplyTextAndAt(fmt.Sprintf("【%s】:\n%s", word, data))
			}
		} else {
			ctx.ReplyText("查询失败，这一定不是bug🤔")
		}
	})
}

func transPinYinSuoXie(text string) (string, error) {
	api := "https://lab.magiconch.com/api/nbnhhsh/guess"
	resp := req.C().Post(api).SetFormData(map[string]string{"text": text}).Do()
	var ret []string
	gjson.Get(resp.String(), "0.trans").ForEach(func(key, val gjson.Result) bool {
		ret = append(ret, val.String())
		return true
	})
	return strings.Join(ret, "；"), nil
}
