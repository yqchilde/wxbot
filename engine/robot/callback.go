package robot

const (
	EventGroupChat           = "EventGroupChat"           // 群聊消息事件
	EventPrivateChat         = "EventPrivateChat"         // 私聊消息事件
	EventMPChat              = "EventMPChat"              // 公众号消息事件
	EventSelfMessage         = "EventSelfMessage"         // 自己发的消息事件
	EventFriendVerify        = "EventFriendVerify"        // 好友请求事件
	EventTransfer            = "EventTransfer"            // 好友转账事件
	EventMessageWithdraw     = "EventMessageWithdraw"     // 消息撤回事件
	EventSystem              = "EventSystem"              // 系统消息事件
	EventGroupMemberIncrease = "EventGroupMemberIncrease" // 群成员增加事件
	EventGroupMemberDecrease = "EventGroupMemberDecrease" // 群成员减少事件
	EventInvitedInGroup      = "EventInvitedInGroup"      // 被邀请入群事件
)

// IsText 判断消息类型是否为文本
func (ctx *Ctx) IsText() bool {
	return ctx.Event.Message != nil && ctx.Event.Message.Type == MsgTypeText
}

// IsImage 判断消息类型是否为图片
func (ctx *Ctx) IsImage() bool {
	return ctx.Event.Message != nil && ctx.Event.Message.Type == MsgTypeImage
}

// IsVoice 判断消息类型是否为语音
func (ctx *Ctx) IsVoice() bool {
	return ctx.Event.Message != nil && ctx.Event.Message.Type == MsgTypeVoice
}

// IsAuthentication 判断消息类型是否是认证消息
func (ctx *Ctx) IsAuthentication() bool {
	return ctx.Event.Message != nil && ctx.Event.Message.Type == MsgTypeAuthentication
}

// IsPossibleFriend 判断消息类型是否是好友推荐消息
func (ctx *Ctx) IsPossibleFriend() bool {
	return ctx.Event.Message != nil && ctx.Event.Message.Type == MsgTypePossibleFriend
}

// IsShareCard 判断消息类型是否是名片消息
func (ctx *Ctx) IsShareCard() bool {
	return ctx.Event.Message != nil && ctx.Event.Message.Type == MsgTypeShareCard
}

// IsVideo 判断消息类型是否是视频消息
func (ctx *Ctx) IsVideo() bool {
	return ctx.Event.Message != nil && ctx.Event.Message.Type == MsgTypeVideo
}

// IsMemePictures 判断消息类型是否为表情包
func (ctx *Ctx) IsMemePictures() bool {
	return ctx.Event.Message != nil && ctx.Event.Message.Type == MsgTypeMemePicture
}

// GetMemePictures 获取表情包图片地址
func (ctx *Ctx) GetMemePictures() (string, bool) {
	if ctx.Event.Message != nil && ctx.Event.Message.Type != MsgTypeMemePicture {
		return "", false
	}
	return ctx.framework.GetMemePictures(ctx.Event.Message), true
}

// IsLocation 判断消息类型是否是地理位置消息
func (ctx *Ctx) IsLocation() bool {
	return ctx.Event.Message != nil && ctx.Event.Message.Type == MsgTypeLocation
}

// IsApp 判断消息类型是否是APP消息
func (ctx *Ctx) IsApp() bool {
	return ctx.Event.Message != nil && ctx.Event.Message.Type == MsgTypeApp
}

// IsMicroVideo 判断消息类型是否是小视频消息
func (ctx *Ctx) IsMicroVideo() bool {
	return ctx.Event.Message != nil && ctx.Event.Message.Type == MsgTypeMicroVideo
}

// IsSystem 判断消息类型是否是系统消息
func (ctx *Ctx) IsSystem() bool {
	return ctx.Event.Message != nil && ctx.Event.Message.Type == MsgTypeSystem
}

// IsRecalled 判断消息类型是否是消息撤回
func (ctx *Ctx) IsRecalled() bool {
	return ctx.Event.Message != nil && ctx.Event.Message.Type == MsgTypeRecalled
}

// IsReference 判断消息类型是否是消息引用
func (ctx *Ctx) IsReference() bool {
	return ctx.Event.Message != nil && ctx.Event.ReferenceMessage != nil
}

// IsAt 判断是否被@了，仅在群聊中有效，私聊也算被@了
func (ctx *Ctx) IsAt() bool {
	return ctx.Event.IsAtMe
}

// IsEventPrivateChat 判断消息是否是私聊消息
func (ctx *Ctx) IsEventPrivateChat() bool {
	return ctx.Event.Type == EventPrivateChat
}

// IsEventGroupChat 判断消息是否是群聊消息
func (ctx *Ctx) IsEventGroupChat() bool {
	return ctx.Event.Type == EventGroupChat
}

// IsEventSelfMessage 判断消息是否是机器人自己发出的消息
func (ctx *Ctx) IsEventSelfMessage() bool {
	return ctx.Event.Type == EventSelfMessage
}

// IsEventFriendVerify 判断消息是否是好友请求消息
func (ctx *Ctx) IsEventFriendVerify() bool {
	return ctx.Event.Type == EventFriendVerify
}

// IsEventSubscription 判断消息是否是订阅消息
func (ctx *Ctx) IsEventSubscription() bool {
	return ctx.Event.Type == EventMPChat
}
