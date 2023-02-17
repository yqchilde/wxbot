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

	// 批量好友发送消息
	engine.OnFullMatch("测试批量好友发送").SetBlock(true).Handle(func(ctx *robot.Ctx) {
		// 写法一
		// 参数为true时从微信数据库拉取数据，false时从启动项目时存入缓存中的数据拉取
		friends, err := ctx.GetFriends()
		if err != nil {
			panic(err)
		}
		friends.SendText("hello", time.Second)

		// 写法二
		// 参数为true时从微信数据库拉取数据，false时从启动项目时存入缓存中的数据拉取
		friends, err = ctx.GetFriends()
		if err != nil {
			panic(err)
		}
		for _, friend := range friends {
			friend.SendText("hello")
			time.Sleep(1 * time.Second)
		}
	})

	// 批量群发送消息
	engine.OnFullMatch("测试批量群发送").SetBlock(true).Handle(func(ctx *robot.Ctx) {
		// 参数为true时从微信数据库拉取数据，false时从启动项目时存入缓存中的数据拉取
		groups, err := ctx.GetGroups()
		if err != nil {
			panic(err)
		}
		for _, group := range groups {
			group.SendText("hello")
			time.Sleep(1 * time.Second)
		}
	})
}
