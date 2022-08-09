package robot

import (
	"sync/atomic"

	"github.com/eatmoreapple/openwechat"

	"wxBot/internal/message"
	"wxBot/internal/model"
)

func Init() {
	// 使用桌面方式登录
	bot := openwechat.DefaultBot(openwechat.Desktop)

	// 关闭心跳回调
	bot.SyncCheckCallback = nil

	// 登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 开启热登录
	reloadStorage := &JsonLocalStorage{FileName: "storage.json"}
	if err := bot.HotLogin(reloadStorage, true); err != nil {
		panic(err)
	}

	// 处理消息回调
	message.NewMessage()
	bot.MessageHandler = func(msg *openwechat.Message) { message.HandleMessage(msg) }

	var count int32
	bot.MessageErrorHandler = func(err error) bool {
		atomic.AddInt32(&count, 1)
		if count == 3 {
			bot.Logout()

		}
		return true
	}

	// 获取登陆的用户
	if self, err := bot.GetCurrentUser(); err == nil {
		model.Self = self
	} else {
		panic(err)
	}

	// 获取所有的好友
	if friends, err := model.Self.Friends(true); err != nil {
		panic(err)
	} else {
		model.Friends = friends
	}

	// 获取所有的群组
	if groups, err := model.Self.Groups(true); err != nil {
		panic(err)
	} else {
		model.Groups = groups
	}

	bot.Block()
}
