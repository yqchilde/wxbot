package crazykfc

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/eatmoreapple/openwechat"
	"github.com/yqchilde/wxbot/engine"
)

type CrazyKFC struct{}

var _ = engine.InstallPlugin(&CrazyKFC{})

func (p *CrazyKFC) OnRegister(event any) {}

func (p *CrazyKFC) OnEvent(event any) {
	if event != nil {
		msg := event.(*openwechat.Message)
		if msg.IsText() && strings.HasPrefix(msg.Content, "/kfc") {
			if resp, err := getCrazyKFCSentence(); err == nil {
				rand.Seed(time.Now().UnixNano())
				msg.ReplyText(resp[rand.Intn(len(resp))].Text)
			} else {
				msg.ReplyText("查询失败，这一定不是bug🤔")
			}
		}
	}
}

type apiResponse struct {
	Index int    `json:"index"`
	Text  string `json:"text"`
}

func getCrazyKFCSentence() ([]apiResponse, error) {
	api := "https://fastly.jsdelivr.net/gh/Nthily/KFC-Crazy-Thursday@main/kfc.json"
	resp, err := http.Get(api)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	readAll, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var data []apiResponse
	if err := json.Unmarshal(readAll, &data); err != nil {
		return nil, err
	}
	return data, nil
}
