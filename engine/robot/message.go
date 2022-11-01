package robot

import (
	"errors"
	"regexp"
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
	if m.IsSendByGroupChat() {
		return MyRobot.SendText(m.Content.FromGroup, msg)
	} else {
		return MyRobot.SendText(m.Content.FromWxid, msg)
	}
}

// ReplyTextAndAt 回复文本消息并@发送者
func (m *Message) ReplyTextAndAt(msg string) error {
	if !m.IsSendByGroupChat() {
		return errors.New("only group chat can reply text and at")
	}
	return MyRobot.SendTextAndAt(msg, m.Content.FromGroup, m.Content.FromWxid, m.Content.FromGroupName)
}

// ReplyImage 回复图片消息
func (m *Message) ReplyImage(path string) error {
	if m.IsSendByGroupChat() {
		return MyRobot.SendImage(m.Content.FromGroup, path)
	} else {
		return MyRobot.SendImage(m.Content.FromWxid, path)
	}
}

// ReplyFile 回复文件消息
func (m *Message) ReplyFile(path string) error {
	if m.IsSendByGroupChat() {
		return MyRobot.SendFile(m.Content.FromGroup, path)
	} else {
		return MyRobot.SendFile(m.Content.FromWxid, path)
	}
}

// ReplyShareLink 回复分享链接消息
func (m *Message) ReplyShareLink(title, desc, imageUrl, jumpUrl string) error {
	if m.IsSendByGroupChat() {
		return MyRobot.SendShareLink(m.Content.FromGroup, title, desc, imageUrl, jumpUrl)
	} else {
		return MyRobot.SendShareLink(m.Content.FromWxid, title, desc, imageUrl, jumpUrl)
	}
}

// WithdrawOwnMessage 撤回自己的消息
func (m *Message) WithdrawOwnMessage() error {
	if m.IsSendByGroupChat() {
		return MyRobot.WithdrawOwnMessage(m.Content.FromGroup, m.Content.MsgId)
	} else {
		return MyRobot.WithdrawOwnMessage(m.Content.FromWxid, m.Content.MsgId)
	}
}

// ReplyVideo 回复视频消息
func (m *Message) ReplyVideo(path string) error {
	if m.IsSendByGroupChat() {
		return MyRobot.SendVideo(m.Content.FromGroup, path)
	} else {
		return MyRobot.SendVideo(m.Content.FromWxid, path)
	}
}
