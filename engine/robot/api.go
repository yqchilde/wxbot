package robot

import "errors"

// IFramework 这是接入框架所定义的接口
type IFramework interface {
	// Callback 这是消息回调方法，vx框架回调消息转发给该Server
	Callback(func(*Event, IFramework))

	// GetMemePictures 判断是否是表情包图片(迷因图)
	// return: 图片链接(网络URL或图片base64)
	GetMemePictures(message *Message) string

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

	// AgreeFriendVerify 同意好友验证
	// v3: 验证V3
	// v4: 验证V4
	// scene: 验证场景
	AgreeFriendVerify(v3, v4, scene string) error

	// InviteIntoGroup 邀请好友加入群组
	// groupWxId: 群ID
	// wxId: 好友ID
	// typ: 邀请类型，1-直接拉，2-发送邀请链接
	InviteIntoGroup(groupWxId, wxId string, typ int) error
}

// SendText 发送文本消息到指定好友
func (ctx *Ctx) SendText(wxId, text string) error {
	if text == "" {
		return nil
	}
	return ctx.framework.SendText(wxId, text)
}

// SendTextAndAt 发送文本消息并@某人到指定群指定用户，仅限群聊
func (ctx *Ctx) SendTextAndAt(groupWxId, wxId, text string) error {
	return ctx.framework.SendTextAndAt(groupWxId, wxId, "", text)
}

// SendImage 发送图片消息到指定好友
func (ctx *Ctx) SendImage(wxId, path string) error {
	return ctx.framework.SendImage(wxId, path)
}

// SendShareLink 发送分享链接消息到指定好友
func (ctx *Ctx) SendShareLink(wxId, title, desc, imageUrl, jumpUrl string) error {
	return ctx.framework.SendShareLink(wxId, title, desc, imageUrl, jumpUrl)
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
	if ctx.IsEventPrivateChat() {
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

// AgreeFriendVerify 同意好友验证
func (ctx *Ctx) AgreeFriendVerify(v3, v4, scene string) error {
	return ctx.framework.AgreeFriendVerify(v3, v4, scene)
}

// InviteIntoGroup 邀请好友加入群组; typ:1-直接拉，2-发送邀请链接
func (ctx *Ctx) InviteIntoGroup(groupWxId, wxId string, typ int) error {
	if typ != 1 && typ != 2 {
		return errors.New("类型错误，请参考方法注释")
	}
	return ctx.framework.InviteIntoGroup(groupWxId, wxId, typ)
}
