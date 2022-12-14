package vlw

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/imroc/req/v3"
	"github.com/yqchilde/pkgs/log"

	"github.com/yqchilde/wxbot/engine/robot"
)

type MessageResp struct {
	Code      int    `json:"Code"`
	Result    string `json:"Result"`
	ReturnStr string `json:"ReturnStr"`
	ReturnInt string `json:"ReturnInt"`
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
			buff.WriteString(`[emoji=\u`)
			buff.WriteString(fmt.Sprintf("%04x", r) + `]`)
		case 4:
			r1, r2 := utf16.EncodeRune(r)
			buff.WriteString(`[emoji=\u`)
			buff.WriteString(strconv.FormatInt(int64(r1), 16))
			buff.WriteString(`\u`)
			buff.WriteString(strconv.FormatInt(int64(r2), 16) + `]`)
		default:
			buff.WriteString(string(r))
		}
	}
	return buff.String()
}

func (f *Framework) GetMemePictures(msg *robot.Message) string {
	// 获取图片base64
	path := msg.Content[5 : len(msg.Content)-1]
	payload := map[string]interface{}{
		"api":   "GetFileFoBase64",
		"token": f.ApiToken,
		"path":  path,
	}

	var resp MessageResp
	err := req.C().Post(f.ApiUrl).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("[VLW] GetFileFoBase64 error: %v", err)
		return ""
	}
	if resp.Code != 0 {
		log.Errorf("[VLW] GetFileFoBase64 error: %s", resp.Result)
		return ""
	}
	return resp.ReturnStr
}

func (f *Framework) SendText(toWxId, text string) error {
	payload := map[string]interface{}{
		"api":        "SendTextMsg",
		"token":      f.ApiToken,
		"msg":        f.msgFormat(text),
		"to_wxid":    toWxId,
		"robot_wxid": f.BotWxId,
	}

	var resp MessageResp
	err := req.C().Post(f.ApiUrl).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("[VLW] SendText error: %v", err)
		return err
	}
	if resp.Code != 0 {
		log.Errorf("[VLW] SendText error: %s", resp.Result)
		return err
	}
	return nil
}

func (f *Framework) SendTextAndAt(toGroupWxId, toWxId, toWxName, text string) error {
	payload := map[string]interface{}{
		"api":         "SendGroupMsgAndAt",
		"token":       f.ApiToken,
		"msg":         f.msgFormat(text),
		"robot_wxid":  f.BotWxId,
		"group_wxid":  toGroupWxId,
		"member_wxid": toWxId,
		"member_name": toWxName,
	}

	var resp MessageResp
	err := req.C().Post(f.ApiUrl).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("[VLW] SendTextAndAt error: %s", err)
		return err
	}
	if resp.Code != 0 {
		log.Errorf("[VLW] SendTextAndAt error: %s", resp.Result)
		return err
	}
	return nil
}

func (f *Framework) SendImage(toWxId, path string) error {
	payload := map[string]interface{}{
		"api":        "SendImageMsg",
		"token":      f.ApiToken,
		"path":       path,
		"robot_wxid": f.BotWxId,
		"to_wxid":    toWxId,
	}

	var resp MessageResp
	err := req.C().Post(f.ApiUrl).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("[VLW] SendImage error: %v", err)
		return err
	}
	if resp.Code != 0 {
		log.Errorf("[VLW] SendImage error: %s", resp.Result)
		return err
	}
	return nil
}

func (f *Framework) SendShareLink(toWxId, title, desc, imageUrl, jumpUrl string) error {
	payload := map[string]interface{}{
		"api":        "SendShareLinkMsg",
		"token":      f.ApiToken,
		"robot_wxid": f.BotWxId,
		"to_wxid":    toWxId,
		"title":      title,
		"desc":       desc,
		"image_url":  imageUrl,
		"url":        jumpUrl,
	}

	var resp MessageResp
	err := req.C().Post(f.ApiUrl).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("[VLW] SendShareLink error: %v", err)
		return err
	}
	if resp.Code != 0 {
		log.Errorf("[VLW] SendShareLink error: %s", resp.Result)
		return err
	}
	return nil
}

func (f *Framework) AgreeFriendVerify(v1, v2, scene string) error {
	// todo 抽空补充
	return nil
}

func (f *Framework) InviteIntoGroup(groupWxId, wxId string, typ int) error {
	// todo 抽空补充
	return nil
}
