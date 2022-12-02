package crazykfc

import (
	"math/rand"
	"time"

	"github.com/imroc/req/v3"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/robot"
)

var sentence []string

func init() {
	resp, err := getCrazyKFCSentence()
	if err != nil {
		return
	}
	for i := range resp {
		sentence = append(sentence, resp[i].Text)
	}

	engine := control.Register("kfccrazy", &control.Options[*robot.Ctx]{
		Alias: "kfc骚话",
		Help:  "输入 {kfc骚话} => 获取肯德基疯狂星期四骚话",
	})

	engine.OnFullMatch("kfc骚话").SetBlock(true).Handle(func(ctx *robot.Ctx) {
		if len(sentence) > 0 {
			rand.Seed(time.Now().UnixNano())
			idx := rand.Intn(len(sentence) - 1)
			ctx.ReplyText(sentence[idx])
			sentence = append(sentence[:idx], sentence[idx+1:]...)
		} else {
			ctx.ReplyText("查询失败，这一定不是bug🤔")
		}
	})
}

type apiResponse struct {
	Index int    `json:"index"`
	Text  string `json:"text"`
}

func getCrazyKFCSentence() ([]apiResponse, error) {
	var data []apiResponse
	api := "https://fastly.jsdelivr.net/gh/Nthily/KFC-Crazy-Thursday@main/kfc.json"
	if err := req.C().Get(api).Do().Into(&data); err != nil {
		return nil, err
	}
	return data, nil
}
