package vlw

import (
	"bytes"
	"encoding/json"
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

	var dataResp MessageResp
	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).SetResult(&dataResp).Do().Err; err != nil {
		log.Errorf("[VLW] GetMemePictures error: %v", err.Error())
		return ""
	}
	return dataResp.ReturnStr
}

func (f *Framework) SendText(toWxId, text string) error {
	payload := map[string]interface{}{
		"api":        "SendTextMsg",
		"token":      f.ApiToken,
		"msg":        f.msgFormat(text),
		"to_wxid":    toWxId,
		"robot_wxid": f.BotWxId,
	}

	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[VLW] SendText error: %v", err.Error())
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

	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[VLW] SendTextAndAt error: %v", err.Error())
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

	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[VLW] SendImage error: %v", err.Error())
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

	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[VLW] SendShareLink error: %v", err.Error())
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

	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[VLW] SendFile error: %v", err.Error())
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

	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[VLW] SendVideo error: %v", err.Error())
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

	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[VLW] SendEmoji error: %v", err.Error())
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

	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[VLW] SendMusic error: %v", err.Error())
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

	buf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(buf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(payload)
	if err := req.C().Post(f.ApiUrl).SetBodyJsonString(buf.String()).Do().Err; err != nil {
		log.Errorf("[VLW] SendMessageRecordXML error: %v", err.Error())
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

	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[VLW] SendFavorites error: %v", err.Error())
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

	buf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(buf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(payload)
	if err := req.C().Post(f.ApiUrl).SetBodyJsonString(buf.String()).Do().Err; err != nil {
		log.Errorf("[VLW] SendXML error: %v", err.Error())
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

	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[VLW] SendBusinessCard error: %v", err.Error())
		return err
	}
	return nil
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

	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[VLW] AgreeFriendVerify error: %v", err.Error())
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
	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[VLW] InviteIntoGroup error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) GetObjectInfo(wxId string) (*robot.ObjectInfo, error) {
	return nil, nil
}
