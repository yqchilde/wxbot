package vlw

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/imroc/req/v3"

	"github.com/yqchilde/wxbot/engine/pkg/log"
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

func (f *Framework) SendFile(toWxId, path string) error {
	payload := map[string]interface{}{
		"api":        "SendFileMsg",
		"token":      f.ApiToken,
		"robot_wxid": f.BotWxId,
		"to_wxid":    toWxId,
		"path":       path,
	}

	var resp MessageResp
	err := req.C().Post(f.ApiUrl).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("[VLW] SendFile error: %v", err)
		return err
	}
	if resp.Code != 0 {
		log.Errorf("[VLW] SendFile error: %s", resp.Result)
		return err
	}
	return nil
}

func (f *Framework) SendVideo(toWxId, path string) error {
	payload := map[string]interface{}{
		"api":        "SendVideoMsg",
		"token":      f.ApiToken,
		"robot_wxid": f.BotWxId,
		"to_wxid":    toWxId,
		"path":       path,
	}

	var resp MessageResp
	err := req.C().Post(f.ApiUrl).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("[VLW] SendVideo error: %v", err)
		return err
	}
	if resp.Code != 0 {
		log.Errorf("[VLW] SendVideo error: %s", resp.Result)
		return err
	}
	return nil
}

func (f *Framework) SendEmoji(toWxId, path string) error {
	payload := map[string]interface{}{
		"api":        "SendEmojiMsg",
		"token":      f.ApiToken,
		"robot_wxid": f.BotWxId,
		"to_wxid":    toWxId,
		"path":       path,
	}

	var resp MessageResp
	err := req.C().Post(f.ApiUrl).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("[VLW] SendEmoji error: %v", err)
		return err
	}
	if resp.Code != 0 {
		log.Errorf("[VLW] SendEmoji error: %s", resp.Result)
		return err
	}
	return nil
}

func (f *Framework) SendMusic(toWxId, name, author, app, jumpUrl, musicUrl, coverUrl string) error {
	payload := map[string]interface{}{
		"api":        "SendMusicLinkMsg",
		"token":      f.ApiToken,
		"robot_wxid": f.BotWxId,
		"to_wxid":    toWxId,
		"title":      name,
		"desc":       author,
		"url":        jumpUrl,
		"dataurl":    musicUrl,
		"thumburl":   coverUrl,
	}

	var resp MessageResp
	err := req.C().Post(f.ApiUrl).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("[VLW] SendMusic error: %v", err)
		return err
	}
	if resp.Code != 0 {
		log.Errorf("[VLW] SendMusic error: %s", resp.Result)
		return err
	}
	return nil
}

func (f *Framework) SendMiniProgram(toWxId, ghId, title, content, imagePath, jumpPath string) error {
	log.Errorf("[VLW] SendMiniProgram not support")
	return errors.New("SendMiniProgram not support，please use SendXML")
}

func (f *Framework) SendMessageRecord(toWxId, title string, dataList []map[string]interface{}) error {
	log.Errorf("[千寻] SendMessageRecord not support")
	return errors.New("SendMessageRecord not support, please use SendMessageRecordXML")
}

func (f *Framework) SendMessageRecordXML(toWxId, xmlStr string) error {
	payload := map[string]interface{}{
		"api":        "SendMessageRecord",
		"token":      f.ApiToken,
		"robot_wxid": f.BotWxId,
		"to_wxid":    toWxId,
		"content":    xmlStr,
	}

	var resp MessageResp
	err := req.C().Post(f.ApiUrl).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("[VLW] SendMessageRecordXML error: %v", err)
		return err
	}
	if resp.Code != 0 {
		log.Errorf("[VLW] SendMessageRecordXML error: %s", resp.Result)
		return err
	}
	return nil
}

func (f *Framework) SendFavorites(toWxId, favoritesId string) error {
	payload := map[string]interface{}{
		"api":        "SendFavorites",
		"token":      f.ApiToken,
		"robot_wxid": f.BotWxId,
		"to_wxid":    toWxId,
		"local_id":   favoritesId,
	}

	var resp MessageResp
	err := req.C().Post(f.ApiUrl).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("[VLW] SendFavorites error: %v", err)
		return err
	}
	if resp.Code != 0 {
		log.Errorf("[VLW] SendFavorites error: %s", resp.Result)
		return err
	}
	return nil
}

func (f *Framework) SendXML(toWxId, xmlStr string) error {
	payload := map[string]interface{}{
		"api":        "SendXmlMsg",
		"token":      f.ApiToken,
		"robot_wxid": f.BotWxId,
		"to_wxid":    toWxId,
		"xml":        xmlStr,
	}

	var resp MessageResp
	err := req.C().Post(f.ApiUrl).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("[VLW] SendXML error: %v", err)
		return err
	}
	if resp.Code != 0 {
		log.Errorf("[VLW] SendXML error: %s", resp.Result)
		return err
	}
	return nil
}

func (f *Framework) SendBusinessCard(toWxId, targetWxId string) error {
	payload := map[string]interface{}{
		"api":        "SendCardMsg",
		"token":      f.ApiToken,
		"robot_wxid": f.BotWxId,
		"to_wxid":    toWxId,
		"content":    targetWxId,
	}

	var resp MessageResp
	err := req.C().Post(f.ApiUrl).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("[VLW] SendBusinessCard error: %v", err)
		return err
	}
	if resp.Code != 0 {
		log.Errorf("[VLW] SendBusinessCard error: %s", resp.Result)
		return err
	}
	return nil
}

func (f *Framework) SendBusinessCardXML(toWxId, xmlStr string) error {
	log.Errorf("[千寻] SendBusinessCardXML not support")
	return errors.New("SendBusinessCardXML not support, please use SendBusinessCard")
}

func (f *Framework) AgreeFriendVerify(v1, v2, scene string) error {
	sceneInt, err := strconv.Atoi(scene)
	if err != nil {
		log.Errorf("[VLW] AgreeFriendVerify error: %v", err)
		return err
	}

	payload := map[string]interface{}{
		"api":        "AgreeFriendVerify",
		"token":      f.ApiToken,
		"robot_wxid": f.BotWxId,
		"v1":         v1,
		"v2":         v2,
		"type":       sceneInt,
	}

	var resp MessageResp
	err = req.C().Post(f.ApiUrl).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("[VLW] SendXML error: %v", err)
		return err
	}
	if resp.Code != 0 {
		log.Errorf("[VLW] SendXML error: %s", resp.Result)
		return err
	}
	return nil
}

func (f *Framework) InviteIntoGroup(groupWxId, wxId string, typ int) error {
	payload := make(map[string]interface{}, 5)
	switch typ {
	case 1:
		payload = map[string]interface{}{
			"api":         "InviteInGroup",
			"token":       f.ApiToken,
			"robot_wxid":  f.BotWxId,
			"group_wxid":  groupWxId,
			"friend_wxid": wxId,
		}
	case 2:
		payload = map[string]interface{}{
			"api":         "InviteInGroupByLink",
			"token":       f.ApiToken,
			"robot_wxid":  f.BotWxId,
			"group_wxid":  groupWxId,
			"friend_wxid": wxId,
		}
	}
	var resp MessageResp
	err := req.C().Post(f.ApiUrl).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("[VLW] InviteIntoGroup error: %v", err)
		return err
	}
	if resp.Code != 0 {
		log.Errorf("[VLW] InviteIntoGroup error: %s", resp.Result)
		return err
	}
	return nil
}
