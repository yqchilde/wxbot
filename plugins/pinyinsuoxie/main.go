package pinyinsuoxie

import (
	"bytes"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"

	"github.com/eatmoreapple/openwechat"
	"github.com/tidwall/gjson"

	"github.com/yqchilde/wxbot/engine"
)

type PinYinSuoXie struct{}

var _ = engine.InstallPlugin(&PinYinSuoXie{})

func (p *PinYinSuoXie) OnRegister(event any) {}

func (p *PinYinSuoXie) OnEvent(event any) {
	if event != nil {
		msg := event.(*openwechat.Message)
		if msg.IsText() && strings.HasPrefix(msg.Content, "/??") {
			var re = regexp.MustCompile(`(?m)^/[?？]{1,2} ?([a-z0-9]+)$`)
			match := re.FindAllStringSubmatch(msg.Content, -1)
			if len(match) > 0 && len(match[0]) > 1 {
				if data, err := transPinYinSuoXie(match[0][1]); err == nil {
					msg.ReplyText(match[0][1] + ": " + data)
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

	body, err := ioutil.ReadAll(res.Body)
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
