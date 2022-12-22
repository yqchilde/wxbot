package robot

import (
	"errors"
)

// IFramework 这是接入框架所定义的接口
type IFramework interface {
	// Callback 这是消息回调方法，vx框架回调消息转发给该Server
	Callback(func(*Event, IFramework))

	// GetMemePictures 获取表情包图片地址(迷因图)
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

	// SendFile 发送文件消息
	// toWxId: 好友ID/群ID/公众号ID
	// path: 本地文件绝对路径
	SendFile(toWxId, path string) error

	// SendVideo 发送视频消息
	// toWxId: 好友ID/群ID/公众号ID
	// path: 本地视频文件绝对路径
	SendVideo(toWxId, path string) error

	// SendEmoji 发送表情消息
	// toWxId: 好友ID/群ID/公众号ID
	// path: 本地动态表情文件绝对路径
	SendEmoji(toWxId, path string) error

	// SendMusic 发送音乐消息
	// toWxId: 好友ID/群ID/公众号ID
	// name: 音乐名称
	// author: 音乐作者
	// app: 音乐来源(VLW需留空)，酷狗/wx79f2c4418704b4f8，网易云/wx8dd6ecd81906fd84，QQ音乐/wx5aa333606550dfd5
	// jumpUrl: 音乐跳转链接
	// musicUrl: 网络歌曲直链
	// coverUrl: 封面图片链接
	SendMusic(toWxId, name, author, app, jumpUrl, musicUrl, coverUrl string) error

	// SendMiniProgram 发送小程序消息
	// toWxId: 好友ID/群ID/公众号ID
	// ghId: 小程序ID
	// title: 标题
	// content: 内容
	// imagePath: 图片路径, 本地图片路径或网络图片URL
	// jumpPath: 小程序点击跳转地址，例如：pages/index/index.html
	SendMiniProgram(toWxId, ghId, title, content, imagePath, jumpPath string) error

	// SendMessageRecord 发送消息记录
	// toWxId: 好友ID/群ID/公众号ID
	// title: 仅供电脑上显示用，手机上的话微信会根据[显示昵称]来自动生成 谁和谁的聊天记录
	// dataList:
	// 	- wxid: 发送此条消息的人的wxid
	// 	- nickName: 显示的昵称(可随意伪造)
	// 	- timestamp: 10位时间戳
	// 	- msg: 消息内容
	SendMessageRecord(toWxId, title string, dataList []map[string]interface{}) error

	// SendMessageRecordXML 发送消息记录(XML方式)
	// toWxId: 好友ID/群ID/公众号ID
	// xmlStr: 消息记录XML代码
	SendMessageRecordXML(toWxId, xmlStr string) error

	// SendFavorites 发送收藏消息
	// toWxId: 好友ID/群ID/公众号ID
	// favoritesId: 收藏夹ID
	SendFavorites(toWxId, favoritesId string) error

	// SendXML 发送XML消息
	// toWxId: 好友ID/群ID/公众号ID
	// xmlStr: XML代码
	SendXML(toWxId, xmlStr string) error

	// SendBusinessCard 发送名片消息
	// toWxId: 好友ID/群ID/公众号ID
	// targetWxId: 目标用户ID
	SendBusinessCard(toWxId, targetWxId string) error

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

	// GetObjectInfo 获取对象信息
	// wxId: 好友ID/群ID/公众号ID
	// return: ObjectInfo, error
	GetObjectInfo(wxId string) (*ObjectInfo, error)
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

// SendFile 发送文件消息到指定好友
func (ctx *Ctx) SendFile(wxId, path string) error {
	return ctx.framework.SendFile(wxId, path)
}

// SendVideo 发送视频消息到指定好友
func (ctx *Ctx) SendVideo(wxId, path string) error {
	return ctx.framework.SendVideo(wxId, path)
}

// SendEmoji 发送表情消息到指定好友
func (ctx *Ctx) SendEmoji(wxId, path string) error {
	return ctx.framework.SendEmoji(wxId, path)
}

// SendMusic 发送音乐消息到指定好友
func (ctx *Ctx) SendMusic(toWxId, name, author, app, jumpUrl, musicUrl, coverUrl string) error {
	return ctx.framework.SendMusic(toWxId, name, author, app, jumpUrl, musicUrl, coverUrl)
}

// SendMiniProgram 发送小程序消息到指定好友
func (ctx *Ctx) SendMiniProgram(toWxId, ghId, title, content, imagePath, jumpPath string) error {
	return ctx.framework.SendMiniProgram(toWxId, ghId, title, content, imagePath, jumpPath)
}

// SendMessageRecord 发送消息记录到指定好友
func (ctx *Ctx) SendMessageRecord(toWxId, title string, dataList []map[string]interface{}) error {
	return ctx.framework.SendMessageRecord(toWxId, title, dataList)
}

// SendMessageRecordXML 发送消息记录(XML方式)到指定好友
func (ctx *Ctx) SendMessageRecordXML(toWxId, xmlStr string) error {
	return ctx.framework.SendMessageRecordXML(toWxId, xmlStr)
}

// SendFavorites 发送收藏消息到指定好友
func (ctx *Ctx) SendFavorites(toWxId, favoritesId string) error {
	return ctx.framework.SendFavorites(toWxId, favoritesId)
}

// SendXML 发送XML消息到指定好友
func (ctx *Ctx) SendXML(toWxId, xmlStr string) error {
	return ctx.framework.SendXML(toWxId, xmlStr)
}

// SendBusinessCard 发送名片消息到指定好友
func (ctx *Ctx) SendBusinessCard(toWxId, targetWxId string) error {
	return ctx.framework.SendBusinessCard(toWxId, targetWxId)
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

// ReplyFile 回复文件消息
func (ctx *Ctx) ReplyFile(path string) error {
	return ctx.framework.SendFile(ctx.Event.FromUniqueID, path)
}

// ReplyVideo 回复视频消息
func (ctx *Ctx) ReplyVideo(path string) error {
	return ctx.framework.SendVideo(ctx.Event.FromUniqueID, path)
}

// ReplyEmoji 回复表情消息
func (ctx *Ctx) ReplyEmoji(path string) error {
	return ctx.framework.SendEmoji(ctx.Event.FromUniqueID, path)
}

// ReplyMusic 回复音乐消息
func (ctx *Ctx) ReplyMusic(name, author, app, jumpUrl, musicUrl, coverUrl string) error {
	return ctx.framework.SendMusic(ctx.Event.FromUniqueID, name, author, app, jumpUrl, musicUrl, coverUrl)
}

// ReplyMiniProgram 回复小程序消息
func (ctx *Ctx) ReplyMiniProgram(ghId, title, content, imagePath, jumpPath string) error {
	return ctx.framework.SendMiniProgram(ctx.Event.FromUniqueID, ghId, title, content, imagePath, jumpPath)
}

// ReplyMessageRecord 回复消息记录
func (ctx *Ctx) ReplyMessageRecord(title string, dataList []map[string]interface{}) error {
	return ctx.framework.SendMessageRecord(ctx.Event.FromUniqueID, title, dataList)
}

// ReplyMessageRecordXML 回复消息记录(XML方式)
func (ctx *Ctx) ReplyMessageRecordXML(xmlStr string) error {
	return ctx.framework.SendMessageRecordXML(ctx.Event.FromUniqueID, xmlStr)
}

// ReplyFavorites 回复收藏消息
func (ctx *Ctx) ReplyFavorites(favoritesId string) error {
	return ctx.framework.SendFavorites(ctx.Event.FromUniqueID, favoritesId)
}

// ReplyXML 回复XML消息
func (ctx *Ctx) ReplyXML(xmlStr string) error {
	return ctx.framework.SendXML(ctx.Event.FromUniqueID, xmlStr)
}

// ReplyBusinessCard 回复名片消息
func (ctx *Ctx) ReplyBusinessCard(targetWxId string) error {
	return ctx.framework.SendBusinessCard(ctx.Event.FromUniqueID, targetWxId)
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
