package qianxun

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/antchfx/xmlquery"
	"github.com/imroc/req/v3"

	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/robot"
)

type MessageResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (f *Framework) msgFormat(msg string) string {
	buff := bytes.NewBuffer(make([]byte, 0, len(msg)*2))
	for _, r := range msg {
		if unicode.Is(unicode.Han, r) || unicode.IsLetter(r) {
			buff.WriteString(string(r))
			continue
		}
		switch utf8.RuneLen(r) {
		case 2, 3:
			buff.WriteString(`[emoji=`)
			buff.WriteString(fmt.Sprintf("%04x", r) + `]`)
		case 4:
			r1, r2 := utf16.EncodeRune(r)
			buff.WriteString(`[emoji=`)
			buff.WriteString(fmt.Sprintf("%04x]", r1))
			buff.WriteString(`[emoji=`)
			buff.WriteString(fmt.Sprintf("%04x]", r2))
		default:
			buff.WriteString(string(r))
		}
	}
	return strings.ReplaceAll(strings.ReplaceAll(buff.String(), "\r\n", "\r"), "\n", "\r")
}

func (f *Framework) GetMemePictures(msg *robot.Message) string {
	doc, err := xmlquery.Parse(strings.NewReader(msg.Content))
	if err != nil {
		return ""
	}
	node, err := xmlquery.Query(doc, "//emoji")
	if err != nil {
		return ""
	}
	return node.SelectAttr("cdnurl")
}

func (f *Framework) SendText(toWxId, text string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/httpapi/?wxid=%s", f.ApiUrl, f.BotWxId)
	payload := map[string]interface{}{
		"type": "Q0001",
		"data": map[string]interface{}{
			"wxid": toWxId,
			"msg":  f.msgFormat(text),
		},
	}

	var msgResp MessageResp
	resp := req.C().Post(apiUrl).SetBody(payload).Do()
	if err := resp.Into(&msgResp); err != nil {
		log.Errorf("[千寻] SendText error: %v", err)
		return err
	}
	if msgResp.Code != 200 {
		log.Errorf("[千寻] SendText error: %s", resp.String())
		return errors.New(msgResp.Msg)
	}
	return nil
}

func (f *Framework) SendTextAndAt(toGroupWxId, toWxId, toWxName, text string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/httpapi/?wxid=%s", f.ApiUrl, f.BotWxId)
	payload := map[string]interface{}{
		"type": "Q0001",
		"data": map[string]interface{}{
			"wxid": toGroupWxId,
			"msg":  fmt.Sprintf("[@,wxid=%s,nick=%s,isAuto=true] %s", toWxId, toWxName, f.msgFormat(text)),
		},
	}

	var msgResp MessageResp
	resp := req.C().Post(apiUrl).SetBody(payload).Do()
	if err := resp.Into(&msgResp); err != nil {
		log.Errorf("[千寻] SendTextAndAt error: %v", err)
		return err
	}
	if msgResp.Code != 200 {
		log.Errorf("[千寻] SendTextAndAt error: %s", resp.String())
		return errors.New(msgResp.Msg)
	}
	return nil
}

func (f *Framework) SendImage(toWxId, path string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/httpapi/?wxid=%s", f.ApiUrl, f.BotWxId)
	payload := map[string]interface{}{
		"type": "Q0010",
		"data": map[string]interface{}{
			"wxid": toWxId,
			"path": path,
		},
	}

	var msgResp MessageResp
	resp := req.C().Post(apiUrl).SetBody(payload).Do()
	if err := resp.Into(&msgResp); err != nil {
		log.Errorf("[千寻] SendImage error: %v", err)
		return err
	}
	if msgResp.Code != 200 {
		log.Errorf("[千寻] SendImage error: %s", resp.String())
		return errors.New(msgResp.Msg)
	}
	return nil
}

func (f *Framework) SendShareLink(toWxId, title, desc, imageUrl, jumpUrl string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/httpapi/?wxid=%s", f.ApiUrl, f.BotWxId)
	payload := map[string]interface{}{
		"type": "Q0012",
		"data": map[string]interface{}{
			"wxid":    toWxId,
			"title":   title,
			"content": desc,
			"jumpUrl": jumpUrl,
			"path":    imageUrl,
		},
	}

	var msgResp MessageResp
	resp := req.C().Post(apiUrl).SetBody(payload).Do()
	if err := resp.Into(&msgResp); err != nil {
		log.Errorf("[千寻] SendShareLink error: %v", err)
		return err
	}
	if msgResp.Code != 200 {
		log.Errorf("[千寻] SendShareLink error: %s", resp.String())
		return errors.New(msgResp.Msg)
	}
	return nil
}

func (f *Framework) SendFile(toWxId, path string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/httpapi/?wxid=%s", f.ApiUrl, f.BotWxId)
	payload := map[string]interface{}{
		"type": "Q0011",
		"data": map[string]interface{}{
			"wxid": toWxId,
			"path": path,
		},
	}

	var msgResp MessageResp
	resp := req.C().Post(apiUrl).SetBody(payload).Do()
	if err := resp.Into(&msgResp); err != nil {
		log.Errorf("[千寻] SendFile error: %v", err)
		return err
	}
	if msgResp.Code != 200 {
		log.Errorf("[千寻] SendFile error: %s", resp.String())
		return errors.New(msgResp.Msg)
	}
	return nil
}

func (f *Framework) SendVideo(toWxId, path string) error {
	log.Errorf("[千寻] SendVideo not support")
	return errors.New("SendVideo not support")
}

func (f *Framework) SendEmoji(toWxId, path string) error {
	log.Errorf("[千寻] SendEmoji not support")
	return errors.New("SendEmoji not support")
}

func (f *Framework) SendMusic(toWxId, name, author, app, jumpUrl, musicUrl, coverUrl string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/httpapi/?wxid=%s", f.ApiUrl, f.BotWxId)
	payload := map[string]interface{}{
		"type": "Q0014",
		"data": map[string]interface{}{
			"wxid":     toWxId,
			"name":     name,
			"author":   author,
			"app":      app,
			"jumpUrl":  jumpUrl,
			"musicUrl": musicUrl,
			"imageUrl": coverUrl,
		},
	}

	var msgResp MessageResp
	resp := req.C().Post(apiUrl).SetBody(payload).Do()
	if err := resp.Into(&msgResp); err != nil {
		log.Errorf("[千寻] SendMusic error: %v", err)
		return err
	}
	if msgResp.Code != 200 {
		log.Errorf("[千寻] SendMusic error: %s", resp.String())
		return errors.New(msgResp.Msg)
	}
	return nil
}

func (f *Framework) SendMiniProgram(toWxId, ghId, title, content, imagePath, jumpPath string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/httpapi/?wxid=%s", f.ApiUrl, f.BotWxId)
	payload := map[string]interface{}{
		"type": "Q0013",
		"data": map[string]interface{}{
			"wxid":     toWxId,
			"title":    title,
			"content":  content,
			"jumpPath": jumpPath,
			"gh":       ghId,
			"path":     imagePath,
		},
	}

	var msgResp MessageResp
	resp := req.C().Post(apiUrl).SetBody(payload).Do()
	if err := resp.Into(&msgResp); err != nil {
		log.Errorf("[千寻] SendMiniProgram error: %v", err)
		return err
	}
	if msgResp.Code != 200 {
		log.Errorf("[千寻] SendMiniProgram error: %s", resp.String())
		return errors.New(msgResp.Msg)
	}
	return nil
}

func (f *Framework) SendMessageRecord(toWxId, title string, dataList []map[string]interface{}) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/httpapi/?wxid=%s", f.ApiUrl, f.BotWxId)
	payload := map[string]interface{}{
		"type": "Q0009",
		"data": map[string]interface{}{
			"wxid":     toWxId,
			"title":    title,
			"dataList": dataList,
		},
	}

	var msgResp MessageResp
	resp := req.C().Post(apiUrl).SetBody(payload).Do()
	if err := resp.Into(&msgResp); err != nil {
		log.Errorf("[千寻] SendMessageRecord error: %v", err)
		return err
	}
	if msgResp.Code != 200 {
		log.Errorf("[千寻] SendMessageRecord error: %s", resp.String())
		return errors.New(msgResp.Msg)
	}
	return nil
}

func (f *Framework) SendMessageRecordXML(toWxId, xmlStr string) error {
	log.Errorf("[千寻] SendMessageRecordXML not support")
	return errors.New("SendMessageRecordXML not support, please use SendMessageRecord")
}

func (f *Framework) SendFavorites(toWxId, favoritesId string) error {
	log.Errorf("[千寻] SendFavorites not support")
	return errors.New("SendFavorites not support")
}

func (f *Framework) SendXML(toWxId, xmlStr string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/httpapi/?wxid=%s", f.ApiUrl, f.BotWxId)
	payload := map[string]interface{}{
		"type": "Q0015",
		"data": map[string]interface{}{
			"wxid": toWxId,
			"xml":  xmlStr,
		},
	}

	var msgResp MessageResp
	resp := req.C().Post(apiUrl).SetBody(payload).Do()
	if err := resp.Into(&msgResp); err != nil {
		log.Errorf("[千寻] SendXML error: %v", err)
		return err
	}
	if msgResp.Code != 200 {
		log.Errorf("[千寻] SendXML error: %s", resp.String())
		return errors.New(msgResp.Msg)
	}
	return nil
}

func (f *Framework) SendBusinessCard(toWxId, targetWxId string) error {
	log.Errorf("[千寻] SendBusinessCard not support")
	return errors.New("SendBusinessCard not support, please use SendBusinessCardXML")
}

func (f *Framework) SendBusinessCardXML(toWxId, xmlStr string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/httpapi/?wxid=%s", f.ApiUrl, f.BotWxId)
	payload := map[string]interface{}{
		"type": "Q0025",
		"data": map[string]interface{}{
			"wxid": toWxId,
			"xml":  xmlStr,
		},
	}

	var msgResp MessageResp
	resp := req.C().Post(apiUrl).SetBody(payload).Do()
	if err := resp.Into(&msgResp); err != nil {
		log.Errorf("[千寻] SendBusinessCardXML error: %v", err)
		return err
	}
	if msgResp.Code != 200 {
		log.Errorf("[千寻] SendBusinessCardXML error: %s", resp.String())
		return errors.New(msgResp.Msg)
	}
	return nil
}

func (f *Framework) AgreeFriendVerify(v3, v4, scene string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/httpapi/?wxid=%s", f.ApiUrl, f.BotWxId)
	payload := map[string]interface{}{
		"type": "Q0017",
		"data": map[string]interface{}{
			"scene": scene,
			"v3":    v3,
			"v4":    v4,
		},
	}

	var msgResp MessageResp
	resp := req.C().Post(apiUrl).SetBody(payload).Do()
	if err := resp.Into(&msgResp); err != nil {
		log.Errorf("[千寻] AgreeFriendVerify error: %v", err)
		return err
	}
	if msgResp.Code != 200 {
		log.Errorf("[千寻] AgreeFriendVerify error: %s", resp.String())
		return errors.New(msgResp.Msg)
	}
	return nil
}

func (f *Framework) InviteIntoGroup(groupWxId, wxId string, typ int) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/httpapi/?wxid=%s", f.ApiUrl, f.BotWxId)
	payload := map[string]interface{}{
		"type": "Q0021",
		"data": map[string]interface{}{
			"wxid":    groupWxId,
			"objWxid": wxId,
			"type":    typ,
		},
	}

	var msgResp MessageResp
	resp := req.C().Post(apiUrl).SetBody(payload).Do()
	if err := resp.Into(&msgResp); err != nil {
		log.Errorf("[千寻] InviteIntoGroup error: %v", err)
		return err
	}
	if msgResp.Code != 200 {
		log.Errorf("[千寻] InviteIntoGroup error: %s", resp.String())
		return errors.New(msgResp.Msg)
	}
	return nil
}
