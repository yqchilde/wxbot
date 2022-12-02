package robot

import (
	"strings"

	"github.com/antchfx/xmlquery"
)

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
	doc, err := xmlquery.Parse(strings.NewReader(ctx.Event.Message.Msg))
	if err != nil {
		return "", false
	}
	node, err := xmlquery.Query(doc, "//emoji")
	if err != nil {
		return "", false
	}
	return node.SelectAttr("cdnurl"), true
}

// ReplyText 回复文本消息
func (ctx *Ctx) ReplyText(text string) error {
	if ctx.IsSendByPrivateChat() {
		return ctx.caller.SendText(ctx.Event.Message.FromWxId, text)
	}
	return ctx.caller.SendText(ctx.Event.Message.FromGroup, text)
}

// ReplyTextAndAt 回复文本消息并@某人，如果在私聊中则不会@某人
func (ctx *Ctx) ReplyTextAndAt(text string) error {
	if ctx.IsSendByPrivateChat() {
		return ctx.caller.SendText(ctx.Event.Message.FromWxId, text)
	}
	return ctx.caller.SendTextAndAt(ctx.Event.Message.FromGroup, ctx.Event.Message.FromWxId, "", text)
}

// ReplyImage 回复图片消息
func (ctx *Ctx) ReplyImage(path string) error {
	if ctx.IsSendByPrivateChat() {
		return ctx.caller.SendImage(ctx.Event.Message.FromWxId, path)
	}
	return ctx.caller.SendImage(ctx.Event.Message.FromGroup, path)
}

// ReplyShareLink 回复分享链接消息
func (ctx *Ctx) ReplyShareLink(title, desc, imageUrl, jumpUrl string) error {
	if ctx.IsSendByPrivateChat() {
		return ctx.caller.SendShareLink(ctx.Event.Message.FromWxId, title, desc, imageUrl, jumpUrl)
	}
	return ctx.caller.SendShareLink(ctx.Event.Message.FromGroup, title, desc, imageUrl, jumpUrl)
}
