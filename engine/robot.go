package engine

import (
	"fmt"
	"sync/atomic"

	"github.com/eatmoreapple/openwechat"
	"github.com/yqchilde/pkgs/log"

	"github.com/yqchilde/wxbot/engine/robot"
)

func InitRobot() {
	// 使用桌面方式登录
	bot := openwechat.DefaultBot(openwechat.Desktop)

	// 关闭心跳回调
	bot.SyncCheckCallback = nil

	// 登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 开启热登录
	reloadStorage := &robot.JsonLocalStorage{FileName: "storage.json"}
	if err := bot.HotLogin(reloadStorage, true); err != nil {
		panic(err)
	}

	// 处理消息回调
	bot.MessageHandler = func(msg *openwechat.Message) {
		if msg.IsSendBySelf() {
			return
		}

		reply := "YY Bot🤖\n"
		for _, plugin := range Plugins {
			if plugin.RawConfig["enable"] != false {
				plugin.Config.OnEvent(msg)
			}
			reply += plugin.Desc + "\n"
		}

		if msg.IsText() && msg.Content == "/menu" {
			msg.ReplyText(reply)
		}
		if msg.IsSendByFriend() {
			sender, err := msg.Sender()
			if err != nil {
				log.Printf("get friend chat sender error: %v", err)
				return
			}

			if msg.IsText() {
				log.Println(fmt.Sprintf("收到私聊(%s)消息 ==> %v", sender.NickName, msg.Content))
			} else {
				log.Println(fmt.Sprintf("收到私聊(%s)消息 ==> %v", sender.NickName, msg.String()))
			}
		} else {
			sender, err := msg.SenderInGroup()
			if err != nil {
				log.Printf("get group chat sender error: %v", err)
				return
			}

			if msg.IsText() {
				log.Println(fmt.Sprintf("收到群(%s[%s])消息 ==> %v", getGroupNicknameByGroupUsername(msg.FromUserName), sender.NickName, msg.Content))
			} else {
				log.Println(fmt.Sprintf("收到群(%s[%s])消息 ==> %v", getGroupNicknameByGroupUsername(msg.FromUserName), sender.NickName, msg.String()))
			}
		}
	}

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
		robot.Self = self
	} else {
		panic(err)
	}

	// 获取所有的好友
	if friends, err := robot.Self.Friends(true); err != nil {
		panic(err)
	} else {
		robot.Friends = friends
	}

	// 获取所有的群组
	if groups, err := robot.Self.Groups(true); err != nil {
		panic(err)
	} else {
		robot.Groups = groups
	}

	bot.Block()
}

func getGroupNicknameByGroupUsername(username string) string {
	groups := robot.Groups.SearchByUserName(1, username)
	return groups[0].NickName
}
