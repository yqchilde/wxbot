package robot

import (
	"errors"
	"regexp"

	"github.com/imroc/req/v3"
	"github.com/yqchilde/pkgs/log"
)

const (
	EventGroupChat           = "EventGroupChat"           // 群聊消息事件
	EventPrivateChat         = "EventPrivateChat"         // 私聊消息事件
	EventDeviceCallback      = "EventDeviceCallback"      // 设备回调事件
	EventFriendVerify        = "EventFrieneVerify"        // 好友请求事件
	EventGroupNameChange     = "EventGroupNameChange"     // 群名称变动事件
	EventGroupMemberAdd      = "EventGroupMemberAdd"      // 群成员增加事件
	EventGroupMemberDecrease = "EventGroupMemberDecrease" // 群成员减少事件
	EventInvitedInGroup      = "EventInvitedInGroup"      // 被邀请入群事件
	EventQRCodePayment       = "EventQRcodePayment"       // 面对面收款事件
	EventDownloadFile        = "EventDownloadFile"        // 文件下载结束事件
	EventGroupEstablish      = "EventGroupEstablish"      // 创建新的群聊事件
)

type Message struct {
	SdkVer  int    `json:"sdkVer"`
	Event   string `json:"Event"`
	Content struct {
		Type          int            `json:"type"`
		RobotWxid     string         `json:"robot_wxid"`
		FromGroup     string         `json:"from_group"`
		FromGroupName string         `json:"from_group_name"`
		FromWxid      string         `json:"from_wxid"`
		FromName      string         `json:"from_name"`
		Msg           string         `json:"msg"`
		MsgSource     *MessageSource `json:"msg_source"`
		Clientid      int            `json:"clientid"`
		RobotType     int            `json:"robot_type"`
		MsgId         string         `json:"msg_id"`
	} `json:"content"`
}

type MessageSource struct {
	Atuserlist []struct {
		Wxid     string `json:"wxid"`
		Nickname string `json:"nickname"`
	} `json:"atuserlist"`
}

type MessageResp struct {
	Code      int    `json:"Code"`
	Result    string `json:"Result"`
	ReturnStr string `json:"ReturnStr"`
}

// IsText 判断消息是否为文本消息
func (m *Message) IsText() bool {
	return m.Content.Type == MsgTypeText
}

// IsEmoticon 判断消息是否为表情消息
func (m *Message) IsEmoticon() bool {
	return m.Content.Type == MsgTypeEmoticon
}

// IsAt 判断是否被@了
func (m *Message) IsAt() bool {
	if !m.IsSendByGroupChat() {
		return false
	}
	if m.Content.MsgSource != nil {
		if len(m.Content.MsgSource.Atuserlist) > 0 && m.Content.MsgSource.Atuserlist[0].Wxid == m.Content.RobotWxid {
			return true
		}
	}
	return false
}

// IsSendByGroupChat 判断消息是否为群聊消息
func (m *Message) IsSendByGroupChat() bool {
	return m.Event == EventGroupChat
}

// IsSendByPrivateChat 判断消息是否为私聊消息
func (m *Message) IsSendByPrivateChat() bool {
	return m.Event == EventPrivateChat
}

// IsDeviceCallback 判断消息是否为设备回调消息
func (m *Message) IsDeviceCallback() bool {
	return m.Event == EventDeviceCallback
}

// MatchTextCommand 判断消息是否为指定的文本命令
func (m *Message) MatchTextCommand(commands []string) (match bool) {
	if m.IsText() {
		for i := range commands {
			if commands[i] == m.Content.Msg {
				return true
			}
		}
	}
	return false
}

// MatchRegexCommand 判断消息是否为指定的正则命令
func (m *Message) MatchRegexCommand(commands []string) (index int, match bool) {
	if m.IsText() {
		for i := range commands {
			re := regexp.MustCompile(commands[i])
			if re.MatchString(m.Content.Msg) {
				return i, true
			}
		}
	}
	return 0, false
}

// ReplyText 回复文本消息
func (m *Message) ReplyText(msg string) error {
	payload := map[string]interface{}{
		"api":        "SendTextMsg",
		"token":      MyRobot.Token,
		"msg":        formatTextMessage(msg),
		"robot_wxid": m.Content.RobotWxid,
		"to_wxid":    m.Content.FromWxid,
	}
	if m.IsSendByGroupChat() {
		payload["to_wxid"] = m.Content.FromGroup
	}

	var resp MessageResp
	err := req.C().Post(MyRobot.Server).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("reply text message error: %v", err)
		return err
	}
	if resp.Code != 0 {
		log.Errorf("reply text message error: %s", resp.Result)
		return err
	}
	return nil
}

// ReplyTextAndAt 回复文本消息并@发送者
func (m *Message) ReplyTextAndAt(msg string) error {
	if !m.IsSendByGroupChat() {
		return errors.New("only group chat can reply text and at")
	}
	payload := map[string]interface{}{
		"api":         "SendGroupMsgAndAt",
		"token":       MyRobot.Token,
		"msg":         formatTextMessage(msg),
		"robot_wxid":  m.Content.RobotWxid,
		"group_wxid":  m.Content.FromGroup,
		"member_wxid": m.Content.FromWxid,
		"member_name": m.Content.FromName,
	}

	var resp MessageResp
	err := req.C().Post(MyRobot.Server).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("reply text message error: %v", err)
		return err
	}
	if resp.Code != 0 {
		log.Errorf("reply text message error: %s", resp.Result)
		return err
	}
	return nil
}

// ReplyImage 回复图片消息
func (m *Message) ReplyImage(path string) error {
	payload := map[string]interface{}{
		"api":        "SendImageMsg",
		"token":      MyRobot.Token,
		"path":       path,
		"robot_wxid": m.Content.RobotWxid,
		"to_wxid":    m.Content.FromWxid,
	}
	if m.IsSendByGroupChat() {
		payload["to_wxid"] = m.Content.FromGroup
	}

	var resp MessageResp
	err := req.C().Post(MyRobot.Server).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("reply image message error: %v", err)
		return err
	}
	if resp.Code != 0 {
		log.Errorf("reply image message error: %s", resp.Result)
		return err
	}
	return nil
}

// ReplyFile 回复文件消息
func (m *Message) ReplyFile(path string) error {
	payload := map[string]interface{}{
		"api":        "SendFileMsg",
		"token":      MyRobot.Token,
		"path":       path,
		"robot_wxid": m.Content.RobotWxid,
		"to_wxid":    m.Content.FromWxid,
	}
	if m.IsSendByGroupChat() {
		payload["to_wxid"] = m.Content.FromGroup
	}

	var resp MessageResp
	err := req.C().Post(MyRobot.Server).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("reply file message error: %v", err)
		return err
	}
	if resp.Code != 0 {
		log.Errorf("reply file message error: %s", resp.Result)
		return err
	}
	return nil
}

// ReplyShareLink 回复分享链接消息
func (m *Message) ReplyShareLink(title, desc, imageUrl, jumpUrl string) error {
	payload := map[string]interface{}{
		"api":        "SendShareLinkMsg",
		"token":      MyRobot.Token,
		"robot_wxid": m.Content.RobotWxid,
		"to_wxid":    m.Content.FromWxid,
		"title":      title,
		"desc":       desc,
		"image_url":  imageUrl,
		"url":        jumpUrl,
	}
	if m.IsSendByGroupChat() {
		payload["to_wxid"] = m.Content.FromGroup
	}

	var resp MessageResp
	err := req.C().Post(MyRobot.Server).SetBody(payload).Do().Into(&resp)
	if err != nil {
		log.Errorf("reply share link message error: %v", err)
		return err
	}
	if resp.Code != 0 {
		log.Errorf("reply share link message error: %s", resp.Result)
		return err
	}
	return nil
}

// WithdrawOwnMessage 撤回自己的消息
func (m *Message) WithdrawOwnMessage() error {
	payload := map[string]interface{}{
		"api":        "WithdrawOwnMessage",
		"token":      MyRobot.Token,
		"robot_wxid": m.Content.RobotWxid,
		"to_wxid":    m.Content.FromWxid,
		"msgid":      m.Content.MsgId,
	}
	if m.IsSendByGroupChat() {
		payload["to_wxid"] = m.Content.FromGroup
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
