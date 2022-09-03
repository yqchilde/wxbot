package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/eatmoreapple/openwechat"
	"github.com/yqchilde/pkgs/log"
	"github.com/yqchilde/pkgs/timer"

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

		menuItems := "YY Bot🤖\n"
		for _, plugin := range Plugins {
			if plugin.RawConfig["enable"] != false {
				plugin.Config.OnEvent(&robot.Message{Message: msg})
			}
			if !plugin.HiddenMenu {
				menuItems += plugin.Desc + "\n"
			}
		}

		if msg.IsText() {
			// isAt存在bug，需要跟内容才会触发，后续更新
			if msg.IsAt() {
				msg.ReplyText("您可以发送menu | 菜单获取更多姿势😎")
			}
			if msg.Content == "menu" || msg.Content == "菜单" {
				msg.ReplyText(menuItems)
			}
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

	robot.Bot = bot
	go keepalive()
	bot.Block()
}

func keepalive() {
	task := timer.NewTimerTask()
	_, err := task.AddTaskByFunc("keepalive", "0 0/30 * * * *", func() {
		if robot.Bot.Alive() {
			if checkWhetherNeedToLogin() {
				reloadStorage := &robot.JsonLocalStorage{FileName: "storage.json"}
				if err := robot.Bot.HotLogin(reloadStorage, false); err != nil {
					log.Errorf("热登录续命失败, err: %v", err)
					return
				}
				log.Debug("热登录续命成功")
				if err := robot.Bot.DumpHotReloadStorage(); err != nil {
					log.Errorf("热登录数据持久化失败, err: %v", err)
					return
				}
				log.Debug("热登录数据持久化成功")
			}

			helper, err := robot.Self.FileHelper()
			if err != nil {
				log.Errorf("获取文件助手失败, err: %v", err)
				return
			}
			if _, err := helper.SendText(openwechat.ZombieText); err != nil {
				log.Errorf("Robot保活失败, err: %v", err)
				return
			}
			log.Println("Robot保活成功")
		}
	})
	if err != nil {
		log.Errorf("NewScheduled add task error: %v", err)
	}
}

func checkWhetherNeedToLogin() bool {
	storage, err := os.ReadFile("storage.json")
	if err != nil {
		log.Errorf("获取热登录配置失败, err: %v", err)
		return false
	}

	var hotLoginData openwechat.HotReloadStorageItem
	err = json.Unmarshal(storage, &hotLoginData)
	if err != nil {
		log.Errorf("unmarshal hot login storage err: %v", err)
		return false
	}

	for _, cookies := range hotLoginData.Cookies {
		if len(cookies) <= 0 {
			continue
		}

		for _, cookie := range cookies {
			if cookie.Name == "wxsid" {
				gmtLocal, _ := time.LoadLocation("GMT")
				expiresGMTTime, _ := time.ParseInLocation("Mon, 02-Jan-2006 15:04:05 GMT", cookie.RawExpires, gmtLocal)
				expiresLocalTime := expiresGMTTime.In(time.Local)
				overHours := expiresLocalTime.Sub(time.Now().Local()).Hours()
				log.Debugf("距离登录失效还剩%v小时", overHours)
				return overHours < 3
			}
		}
	}
	return false
}

func getGroupNicknameByGroupUsername(username string) string {
	groups := robot.Groups.SearchByUserName(1, username)
	return groups[0].NickName
}
