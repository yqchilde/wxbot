package robot

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/yqchilde/wxbot/engine/pkg/cryptor"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/pkg/utils"
)

// IFramework 这是接入框架所定义的接口
type IFramework interface {
	// Callback 这是消息回调方法，vx框架回调消息转发给该Server
	Callback(*gin.Context, func(*Event, IFramework))

	// GetRobotInfo 获取机器人信息
	// return: User, error
	GetRobotInfo() (*User, error)

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
	// return: User, error
	GetObjectInfo(wxId string) (*User, error)

	// GetFriends 获取好友列表
	// isRefresh: 是否刷新 false-从缓存中获取，true-重新遍历二叉树并刷新缓存
	// return: []*User, error
	GetFriends(isRefresh bool) ([]*User, error)

	// GetGroups 获取群组列表
	// isRefresh: 是否刷新 false-从缓存中获取，true-重新遍历二叉树并刷新缓存
	// return: []*User, error
	GetGroups(isRefresh bool) ([]*User, error)

	// GetGroupMembers 获取群成员列表
	// groupWxId: 群ID
	// isRefresh: 是否刷新 false-从缓存中获取，true-重新遍历二叉树并刷新缓存
	// return: []*User, error
	GetGroupMembers(groupWxId string, isRefresh bool) ([]*User, error)

	// GetMPs 获取公众号订阅列表
	// isRefresh: 是否刷新 false-从缓存中获取，true-重新遍历二叉树并刷新缓存
	// return: []*User, error
	GetMPs(isRefresh bool) ([]*User, error)
}

// SendText 发送文本消息到指定好友
func (ctx *Ctx) SendText(wxId, text string) error {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	if text == "" {
		return nil
	}
	return ctx.framework.SendText(wxId, text)
}

// SendTextAndAt 发送文本消息并@某人到指定群指定用户，仅限群聊
func (ctx *Ctx) SendTextAndAt(groupWxId, wxId, text string) error {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	return ctx.framework.SendTextAndAt(groupWxId, wxId, "", text)
}

// SendTextAndSendEvent 发送文本消息到指定好友并将文本消息压入队列进行插件匹配
func (ctx *Ctx) SendTextAndSendEvent(wxId, text string) error {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	if text == "" {
		return nil
	}
	err := ctx.framework.SendText(wxId, text)
	if err != nil {
		return err
	}

	ctx.SendEvent(wxId, text)
	return nil
}

// SendEvent 将文本消息压入队列进行插件匹配
func (ctx *Ctx) SendEvent(wxId, text string) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	if text == "" {
		return
	}

	// 加入消息监听队列
	event := Event{
		Type:         EventSelfMessage,
		FromUniqueID: wxId,
		Message: &Message{
			Type:    MsgTypeText,
			Content: text,
		},
	}
	eventBuffer.ProcessEvent(&event, ctx.framework)
}

// SendImage 发送图片消息到指定好友
// 支持本地文件，图片路径以local://开头
func (ctx *Ctx) SendImage(wxId, path string) error {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	if strings.HasPrefix(path, "local://") {
		if !utils.CheckPathExists(path[8:]) {
			log.Errorf("[SendImage] 发送图片失败，文件不存在: %s", path[8:])
			return errors.New("发送图片失败，文件不存在")
		}
		if bot.config.ServerAddress == "" {
			log.Errorf("[SendImage] 发送图片失败，请在config.yaml中配置serverAddress项")
			return errors.New("发送图片失败，请在config.yaml中配置serverAddress项")
		}
		filename, err := cryptor.EncryptFilename(fileSecret, path[8:])
		if err != nil {
			log.Errorf("[SendImage] 加密文件名失败: %v", err)
			return err
		}
		path = bot.config.ServerAddress + "/wxbot/static?file=" + filename
	}
	log.Debugf("[SendImage] Path: %s", path)
	return ctx.framework.SendImage(wxId, path)
}

// SendShareLink 发送分享链接消息到指定好友
// 支持本地文件，图片路径以local://开头
func (ctx *Ctx) SendShareLink(wxId, title, desc, imageUrl, jumpUrl string) error {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	if strings.HasPrefix(imageUrl, "local://") {
		if !utils.CheckPathExists(imageUrl[8:]) {
			log.Errorf("[SendShareLink] 发送分享链接失败，文件不存在: %s", imageUrl[8:])
			return errors.New("发送分享链接失败，文件不存在")
		}
		if bot.config.ServerAddress == "" {
			log.Errorf("[SendShareLink] 发送分享链接失败，请在config.yaml中配置serverAddress项")
			return errors.New("发送分享链接失败，请在config.yaml中配置serverAddress项")
		}
		filename, err := cryptor.EncryptFilename(fileSecret, imageUrl[8:])
		if err != nil {
			log.Errorf("[SendShareLink] 加密文件名失败: %v", err)
			return err
		}
		imageUrl = bot.config.ServerAddress + "/wxbot/static?file=" + filename
	}
	log.Debugf("[SendShareLink] imageUrl: %s", imageUrl)
	return ctx.framework.SendShareLink(wxId, title, desc, imageUrl, jumpUrl)
}

// SendFile 发送文件消息到指定好友
func (ctx *Ctx) SendFile(wxId, path string) error {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	return ctx.framework.SendFile(wxId, path)
}

// SendVideo 发送视频消息到指定好友
func (ctx *Ctx) SendVideo(wxId, path string) error {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	return ctx.framework.SendVideo(wxId, path)
}

// SendEmoji 发送表情消息到指定好友
func (ctx *Ctx) SendEmoji(wxId, path string) error {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	return ctx.framework.SendEmoji(wxId, path)
}

// SendMusic 发送音乐消息到指定好友
func (ctx *Ctx) SendMusic(toWxId, name, author, app, jumpUrl, musicUrl, coverUrl string) error {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	return ctx.framework.SendMusic(toWxId, name, author, app, jumpUrl, musicUrl, coverUrl)
}

// SendMiniProgram 发送小程序消息到指定好友
func (ctx *Ctx) SendMiniProgram(toWxId, ghId, title, content, imagePath, jumpPath string) error {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	return ctx.framework.SendMiniProgram(toWxId, ghId, title, content, imagePath, jumpPath)
}

// SendMessageRecord 发送消息记录到指定好友
func (ctx *Ctx) SendMessageRecord(toWxId, title string, dataList []map[string]interface{}) error {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	return ctx.framework.SendMessageRecord(toWxId, title, dataList)
}

// SendMessageRecordXML 发送消息记录(XML方式)到指定好友
func (ctx *Ctx) SendMessageRecordXML(toWxId, xmlStr string) error {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	return ctx.framework.SendMessageRecordXML(toWxId, xmlStr)
}

// SendFavorites 发送收藏消息到指定好友
func (ctx *Ctx) SendFavorites(toWxId, favoritesId string) error {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	return ctx.framework.SendFavorites(toWxId, favoritesId)
}

// SendXML 发送XML消息到指定好友
func (ctx *Ctx) SendXML(toWxId, xmlStr string) error {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	return ctx.framework.SendXML(toWxId, xmlStr)
}

// SendBusinessCard 发送名片消息到指定好友
func (ctx *Ctx) SendBusinessCard(toWxId, targetWxId string) error {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	return ctx.framework.SendBusinessCard(toWxId, targetWxId)
}

// ReplyText 回复文本消息
func (ctx *Ctx) ReplyText(text string) error {
	if text == "" {
		return nil
	}
	return ctx.SendText(ctx.Event.FromUniqueID, text)
}

// ReplyTextAndAt 回复文本消息并@某人，如果在私聊中或自己的消息则不会@
func (ctx *Ctx) ReplyTextAndAt(text string) error {
	if ctx.IsEventPrivateChat() || ctx.IsEventSelfMessage() {
		return ctx.ReplyText(text)
	}
	return ctx.SendTextAndAt(ctx.Event.FromGroup, ctx.Event.FromWxId, text)
}

// ReplyTextAndPushEvent 回复文本消息并将文本消息压入队列进行插件匹配
func (ctx *Ctx) ReplyTextAndPushEvent(text string) error {
	if text == "" {
		return nil
	}
	return ctx.SendTextAndSendEvent(ctx.Event.FromUniqueID, text)
}

// PushEvent 将文本消息压入队列进行插件匹配
func (ctx *Ctx) PushEvent(text string) {
	if text == "" {
		return
	}
	ctx.SendEvent(ctx.Event.FromUniqueID, text)
}

// ReplyImage 回复图片消息
// 支持本地文件，图片路径以local://开头
func (ctx *Ctx) ReplyImage(path string) error {
	return ctx.SendImage(ctx.Event.FromUniqueID, path)
}

// ReplyShareLink 回复分享链接消息
// 支持本地文件，图片路径以local://开头
func (ctx *Ctx) ReplyShareLink(title, desc, imageUrl, jumpUrl string) error {
	return ctx.SendShareLink(ctx.Event.FromUniqueID, title, desc, imageUrl, jumpUrl)
}

// ReplyFile 回复文件消息
func (ctx *Ctx) ReplyFile(path string) error {
	return ctx.SendFile(ctx.Event.FromUniqueID, path)
}

// ReplyVideo 回复视频消息
func (ctx *Ctx) ReplyVideo(path string) error {
	return ctx.SendVideo(ctx.Event.FromUniqueID, path)
}

// ReplyEmoji 回复表情消息
func (ctx *Ctx) ReplyEmoji(path string) error {
	return ctx.SendEmoji(ctx.Event.FromUniqueID, path)
}

// ReplyMusic 回复音乐消息
func (ctx *Ctx) ReplyMusic(name, author, app, jumpUrl, musicUrl, coverUrl string) error {
	return ctx.SendMusic(ctx.Event.FromUniqueID, name, author, app, jumpUrl, musicUrl, coverUrl)
}

// ReplyMiniProgram 回复小程序消息
func (ctx *Ctx) ReplyMiniProgram(ghId, title, content, imagePath, jumpPath string) error {
	return ctx.SendMiniProgram(ctx.Event.FromUniqueID, ghId, title, content, imagePath, jumpPath)
}

// ReplyMessageRecord 回复消息记录
func (ctx *Ctx) ReplyMessageRecord(title string, dataList []map[string]interface{}) error {
	return ctx.SendMessageRecord(ctx.Event.FromUniqueID, title, dataList)
}

// ReplyMessageRecordXML 回复消息记录(XML方式)
func (ctx *Ctx) ReplyMessageRecordXML(xmlStr string) error {
	return ctx.SendMessageRecordXML(ctx.Event.FromUniqueID, xmlStr)
}

// ReplyFavorites 回复收藏消息
func (ctx *Ctx) ReplyFavorites(favoritesId string) error {
	return ctx.SendFavorites(ctx.Event.FromUniqueID, favoritesId)
}

// ReplyXML 回复XML消息
func (ctx *Ctx) ReplyXML(xmlStr string) error {
	return ctx.SendXML(ctx.Event.FromUniqueID, xmlStr)
}

// ReplyBusinessCard 回复名片消息
func (ctx *Ctx) ReplyBusinessCard(targetWxId string) error {
	return ctx.SendBusinessCard(ctx.Event.FromUniqueID, targetWxId)
}

// AgreeFriendVerify 同意好友验证
func (ctx *Ctx) AgreeFriendVerify(v3, v4, scene string) error {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	return ctx.framework.AgreeFriendVerify(v3, v4, scene)
}

// InviteIntoGroup 邀请好友加入群组; typ:1-直接拉，2-发送邀请链接
func (ctx *Ctx) InviteIntoGroup(groupWxId, wxId string, typ int) error {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	if typ != 1 && typ != 2 {
		return errors.New("类型错误，请参考方法注释")
	}
	return ctx.framework.InviteIntoGroup(groupWxId, wxId, typ)
}

// GetRobotInfo 获取机器人信息
func (ctx *Ctx) GetRobotInfo() (*User, error) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	return ctx.framework.GetRobotInfo()
}

// GetObjectInfo 获取对象信息，wxId: 好友ID/群ID/公众号ID
func (ctx *Ctx) GetObjectInfo(wxId string) (*User, error) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	return ctx.framework.GetObjectInfo(wxId)
}

// GetFriends 获取好友列表
func (ctx *Ctx) GetFriends(isRefresh ...bool) (Friends, error) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	return ctx.Bot.self.Friends(isRefresh...)
}

// GetGroups 获取群组列表
func (ctx *Ctx) GetGroups(isRefresh ...bool) (Groups, error) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	return ctx.Bot.self.Groups(isRefresh...)
}

// GetGroupMembers 获取群组成员列表
func (ctx *Ctx) GetGroupMembers(groupWxId string, isRefresh ...bool) (GroupMembers, error) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	return ctx.Bot.self.GroupMembers(groupWxId, isRefresh...)
}

// GetMPs 获取公众号列表
func (ctx *Ctx) GetMPs(isRefresh ...bool) (MPs, error) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	return ctx.Bot.self.MPs(isRefresh...)
}

// FuzzyGetByRemarkOrNick 模糊查询好友、群组、公众号
func (ctx *Ctx) FuzzyGetByRemarkOrNick(remarkOrNick string) []*User {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	var users []*User
	friends := ctx.Bot.self.friends.FuzzyGetByRemarkOrNick(remarkOrNick).AsUsers()
	groups := ctx.Bot.self.groups.FuzzyGetByRemarkOrNick(remarkOrNick).AsUsers()
	mps := ctx.Bot.self.mps.FuzzyGetByRemarkOrNick(remarkOrNick).AsUsers()
	return append(append(append(users, friends...), groups...), mps...)
}
