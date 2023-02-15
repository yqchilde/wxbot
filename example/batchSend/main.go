package batchSend

import (
	"time"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/robot"
)

func init() {
	engine := control.Register("batchSend", &control.Options{
		Alias: "批量发送",
	})

	// 批量好友发送
	engine.OnFullMatch("测试批量好友发送").SetBlock(true).Handle(func(ctx *robot.Ctx) {
		self, err := ctx.Bot.GetSelf()
		if err != nil {
			return
		}
		// 参数为true时从微信数据库拉取数据，false时从启动项目时存入缓存中的数据拉取
		friends, err := self.Friends(true)
		if err != nil {
			return
		}
		for _, friend := range friends {
			friend.SendText("hello")
			// 为了防止被微信封号，注意控制发送间隔
			time.Sleep(1 * time.Second)
		}

		// 从缓存读取便携写法
		friends2 := ctx.Bot.FriendsFromCache()
		for _, friend := range friends2 {
			friend.SendText("hello")
			// 为了防止被微信封号，注意控制发送间隔
			time.Sleep(1 * time.Second)
		}
	})

	// todo 批量群发送
}
