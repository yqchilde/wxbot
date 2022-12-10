package robot

// IFramework 这是接入框架所定义的接口
type IFramework interface {
	// Callback 这是消息回调方法，vx框架回调消息转发给该Server
	Callback(func(*Event, IFramework))

	// GetMemePictures 判断是否是表情包图片(迷因图)
	// return: 图片链接(网络URL或图片base64)
	GetMemePictures(message Message) string

	// SendText 发送文本消息
	// toWxId: 好友ID/群ID
	// text: 文本内容
	SendText(toWxId, text string) error

	// SendTextAndAt 发送文本消息并@，只有群聊有效
	// toGroupWxId: 群ID
	// toWxId: 好友ID/群ID/all
	// toWxName: 好友昵称/群昵称，留空为自动获取
	// text: 文本内容
	SendTextAndAt(toGroupWxId, toWxId, toWxName, text string) error

	// SendImage 发送图片消息
	// toWxId: 好友ID/群ID
	// path: 图片路径
	SendImage(toWxId, path string) error

	// SendShareLink 发送分享链接消息
	// toWxId: 好友ID/群ID
	// title: 标题
	// desc: 描述
	// imageUrl: 图片链接
	// jumpUrl: 跳转链接
	SendShareLink(toWxId, title, desc, imageUrl, jumpUrl string) error
}

// IsSendByPrivateChat 判断消息是否是私聊消息
func (ctx *Ctx) IsSendByPrivateChat() bool {
	return ctx.Event.IsPrivateChat
}

// IsSendByGroupChat 判断消息是否是群聊消息
func (ctx *Ctx) IsSendByGroupChat() bool {
	return ctx.Event.IsGroupChat
}

// IsText 判断消息类型是否为文本
func (ctx *Ctx) IsText() bool {
	return ctx.Event.Message.MsgType == MsgTypeText
}

// IsAt 判断是否被@了，仅在群聊中有效，私聊也算被@了
func (ctx *Ctx) IsAt() bool {
	return ctx.Event.IsAtMe
}

// IsMemePictures 判断消息类型是否为表情包
func (ctx *Ctx) IsMemePictures() (string, bool) {
	if ctx.Event.Message.MsgType != MsgTypeMemePicture {
		return "", false
	}
	return ctx.framework.GetMemePictures(ctx.Event.Message), true
}

// ReplyText 回复文本消息
func (ctx *Ctx) ReplyText(text string) error {
	if text == "" {
		return nil
	}
	return ctx.framework.SendText(ctx.Event.FromUniqueID, text)
}

// ReplyTextAndAt 回复文本消息并@某人，如果在私聊中则不会@某人
func (ctx *Ctx) ReplyTextAndAt(text string) error {
	if ctx.IsSendByPrivateChat() {
		return ctx.framework.SendText(ctx.Event.FromWxId, text)
	}
	return ctx.framework.SendTextAndAt(ctx.Event.FromGroup, ctx.Event.FromWxId, "", text)
}

// ReplyImage 回复图片消息
func (ctx *Ctx) ReplyImage(path string) error {
	return ctx.framework.SendImage(ctx.Event.FromUniqueID, path)
}

// ReplyShareLink 回复分享链接消息
func (ctx *Ctx) ReplyShareLink(title, desc, imageUrl, jumpUrl string) error {
	return ctx.framework.SendShareLink(ctx.Event.FromUniqueID, title, desc, imageUrl, jumpUrl)
}
