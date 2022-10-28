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
	Code   int    `json:"Code"`
	Result string `json:"Result"`
}

func (m *Message) IsText() bool {
	return m.Content.Type == MsgTypeText
}

func (m *Message) IsEmoticon() bool {
	return m.Content.Type == MsgTypeEmoticon
}

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

func (m *Message) IsSendByGroupChat() bool {
	return m.Event == EventGroupChat
}

func (m *Message) IsSendByPrivateChat() bool {
	return m.Event == EventPrivateChat
}

func (m *Message) MatchTextCommand(commands []string) bool {
	if m.IsText() {
		for i := range commands {
			if commands[i] == m.Content.Msg {
				return true
			}
		}
	}
	return false
}

func (m *Message) MatchRegexCommand(commands []string) bool {
	if m.IsText() {
		for i := range commands {
			re := regexp.MustCompile(commands[i])
			return re.MatchString(m.Content.Msg)
		}
	}
	return false
}

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
	err := req.C().SetBaseURL(MyRobot.Server).Post().SetBody(payload).Do().Into(&resp)
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
	err := req.C().SetBaseURL(MyRobot.Server).Post().SetBody(payload).Do().Into(&resp)
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
	err := req.C().SetBaseURL(MyRobot.Server).Post().SetBody(payload).Do().Into(&resp)
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
