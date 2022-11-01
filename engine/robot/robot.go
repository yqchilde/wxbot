package robot

import (
	"github.com/imroc/req/v3"
	"github.com/yqchilde/pkgs/log"
)

var MyRobot BotConf

type BotConf struct {
	Server string `yaml:"server"`
	Token  string `yaml:"token"`
	Bot    Bot    `yaml:"-"`
}

func (b *BotConf) GetRobotInfo() error {
	payload := map[string]interface{}{
		"api":   "GetRobotList",
		"token": MyRobot.Token,
	}

	var resp BotList
	err := req.C().Post(MyRobot.Server).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("get robot info error: %v", err)
		return err
	}
	if resp.Code != 0 {
		log.Errorf("get robot info error: %s", resp.Result)
		return err
	}
	MyRobot.Bot = resp.ReturnJson.Data[0]
	return nil
}

func (b *BotConf) GetGroupList() ([]Group, error) {
	payload := map[string]interface{}{
		"api":        "GetGrouplist",
		"token":      MyRobot.Token,
		"robot_wxid": MyRobot.Bot.Wxid,
		"is_refresh": "0",
	}

	var resp GroupList
	err := req.C().Post(MyRobot.Server).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("get robot info error: %v", err)
		return nil, err
	}
	if resp.Code != 0 {
		log.Errorf("get robot info error: %s", resp.Result)
		return nil, err
	}
	return resp.ReturnJson, nil
}

// SendText 发送文本消息； to_wxid:好友ID/群ID
func (b *BotConf) SendText(toWxId, msg string) error {
	payload := map[string]interface{}{
		"api":        "SendTextMsg",
		"token":      MyRobot.Token,
		"msg":        formatTextMessage(msg),
		"robot_wxid": MyRobot.Bot.Wxid,
		"to_wxid":    toWxId,
	}

	var resp MessageResp
	err := req.C().Post(MyRobot.Server).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("send text message error: %v", err)
		return err
	}
	if resp.Code != 0 {
		log.Errorf("send text message error: %s", resp.Result)
		return err
	}
	return nil
}

// SendTextAndAt 发送文本消息并@； to_wxid:好友ID/群ID
func (b *BotConf) SendTextAndAt(msg, groupWxId, toWxId, toWxName string) error {
	payload := map[string]interface{}{
		"api":         "SendGroupMsgAndAt",
		"token":       MyRobot.Token,
		"msg":         formatTextMessage(msg),
		"robot_wxid":  MyRobot.Bot.Wxid,
		"group_wxid":  groupWxId,
		"member_wxid": toWxId,
		"member_name": toWxName,
	}

	var resp MessageResp
	err := req.C().Post(MyRobot.Server).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("send text message and at error: %s", err)
		return err
	}
	if resp.Code != 0 {
		log.Errorf("send text message and at error: %s", resp.Result)
		return err
	}
	return nil
}

// SendImage 发送图片消息； to_wxid:好友ID/群ID
func (b *BotConf) SendImage(toWxId, path string) error {
	payload := map[string]interface{}{
		"api":        "SendImageMsg",
		"token":      MyRobot.Token,
		"path":       path,
		"robot_wxid": MyRobot.Bot.Wxid,
		"to_wxid":    toWxId,
	}

	var resp MessageResp
	err := req.C().Post(MyRobot.Server).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("send image message error: %v", err)
		return err
	}
	if resp.Code != 0 {
		log.Errorf("send image message error: %s", resp.Result)
		return err
	}
	return nil
}

// SendFile 发送文件消息； to_wxid:好友ID/群ID
func (b *BotConf) SendFile(toWxId, path string) error {
	payload := map[string]interface{}{
		"api":        "SendFileMsg",
		"token":      MyRobot.Token,
		"path":       path,
		"robot_wxid": MyRobot.Bot.Wxid,
		"to_wxid":    toWxId,
	}

	var resp MessageResp
	err := req.C().Post(MyRobot.Server).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("send file message error: %v", err)
		return err
	}
	if resp.Code != 0 {
		log.Errorf("send file message error: %s", resp.Result)
		return err
	}
	return nil
}

// SendShareLink 发送分享链接消息； to_wxid:好友ID/群ID
func (b *BotConf) SendShareLink(toWxId, title, desc, imageUrl, jumpUrl string) error {
	payload := map[string]interface{}{
		"api":        "SendShareLinkMsg",
		"token":      MyRobot.Token,
		"robot_wxid": MyRobot.Bot.Wxid,
		"to_wxid":    toWxId,
		"title":      title,
		"desc":       desc,
		"image_url":  imageUrl,
		"url":        jumpUrl,
	}

	var resp MessageResp
	err := req.C().Post(MyRobot.Server).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("send share link message error: %v", err)
		return err
	}
	if resp.Code != 0 {
		log.Errorf("send share link message error: %s", resp.Result)
		return err
	}
	return nil
}

// WithdrawOwnMessage 撤回自己的消息； to_wxid:好友ID/群ID
func (b *BotConf) WithdrawOwnMessage(toWxId, msgId string) error {
	payload := map[string]interface{}{
		"api":        "WithdrawOwnMessage",
		"token":      MyRobot.Token,
		"robot_wxid": MyRobot.Bot.Wxid,
		"to_wxid":    toWxId,
		"msgid":      msgId,
	}

	var resp MessageResp
	err := req.C().Post(MyRobot.Server).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("withdraw own message error: %v", err)
		return err
	}
	if resp.Code != 0 {
		log.Errorf("withdraw own message error: %s", resp.Result)
		return err
	}
	return nil
}

// SendVideo 发送视频消息； to_wxid:好友ID/群ID
func (b *BotConf) SendVideo(toWxId, path string) error {
	payload := map[string]interface{}{
		"api":        "SendVideoMsg",
		"token":      MyRobot.Token,
		"path":       path,
		"robot_wxid": MyRobot.Bot.Wxid,
		"to_wxid":    toWxId,
	}

	var resp MessageResp
	err := req.C().Post(MyRobot.Server).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("send video message error: %v", err)
		return err
	}
	if resp.Code != 0 {
		log.Errorf("send video message error: %s", resp.Result)
		return err
	}
	return nil
}

func (b *BotConf) GetFileFoBase64(path string) (string, error) {
	payload := map[string]interface{}{
		"api":   "GetFileFoBase64",
		"token": MyRobot.Token,
		"path":  path,
	}

	var resp MessageResp
	err := req.C().Post(MyRobot.Server).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("get file for base64 error: %v", err)
		return "", err
	}
	if resp.Code != 0 {
		log.Errorf("get file for base64 error: %s", resp.Result)
		return "", err
	}
	return resp.ReturnStr, nil
}
