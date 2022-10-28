package pinyinsuoxie

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"

	"github.com/yqchilde/wxbot/engine"
	"github.com/yqchilde/wxbot/engine/robot"
)

type PinYinSuoXie struct{ engine.PluginMagic }

var (
	pluginInfo = &PinYinSuoXie{
		engine.PluginMagic{
			Desc:     "🚀 输入 {查缩写 XX} => 获取拼音缩写翻译，Ps:查缩写 yyds",
			Commands: []string{"^查缩写 ?([a-zA-Z0-9]+)$"},
		},
	}
	_ = engine.InstallPlugin(pluginInfo)
)

func (p *PinYinSuoXie) OnRegister() {}

func (p *PinYinSuoXie) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.MatchRegexCommand(pluginInfo.Commands) {
			var re = regexp.MustCompile(pluginInfo.Commands[0])
			match := re.FindAllStringSubmatch(msg.Content.Msg, -1)
			if len(match) > 0 && len(match[0]) > 1 {
				if data, err := transPinYinSuoXie(match[0][1]); err == nil {
					if len(data) == 0 {
						msg.ReplyText("没查到该缩写含义")
					} else {
						msg.ReplyText(match[0][1] + ": " + data)
					}
				} else {
					msg.ReplyText("查询失败，这一定不是bug🤔")
				}
			}
		}
	}
}

func transPinYinSuoXie(text string) (string, error) {
	url := "https://lab.magiconch.com/api/nbnhhsh/guess"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("text", text)
	err := writer.Close()
	if err != nil {
		log.Println(err)
		return "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Println(err)
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}

	json := gjson.ParseBytes(body)
	ret := make([]string, 0)
	var jsonPath string
	if json.Get("0.trans").Exists() {
		jsonPath = "0.trans"
	} else {
		jsonPath = "0.inputting"
	}
	for _, value := range json.Get(jsonPath).Array() {
		ret = append(ret, value.String())
	}
	return strings.Join(ret, ";"), nil
}
