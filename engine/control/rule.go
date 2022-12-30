package control

import (
	"sync"

	"github.com/yqchilde/wxbot/engine/robot"
)

var (
	once     = sync.Once{}
	managers = NewManager[*robot.Ctx]("data/manager/plugins.db")
)

func newControl(service string, o *Options[*robot.Ctx]) robot.Rule {
	c := managers.NewControl(service, o)
	return func(ctx *robot.Ctx) bool {
		ctx.State["manager"] = c
		return c.Handler(ctx.Event.FromGroup, ctx.Event.FromWxId)
	}
}

func init() {
	once.Do(func() {
		// 启用、禁用某个插件在某个群或某个私聊
		robot.OnCommandGroup([]string{"启用", "禁用"}, robot.UserOrGroupAdmin).SetBlock(true).FirstPriority().Handle(func(ctx *robot.Ctx) {
			args := ctx.State["args"].(string)
			if args == "" {
				return
			}
			service, ok := managers.Lookup(args)
			if !ok {
				ctx.ReplyTextAndAt("没有找到对应插件服务")
				return
			}
			grp := ctx.Event.FromUniqueID
			switch ctx.State["command"].(string) {
			case "启用":
				if err := service.Enable(grp); err != nil {
					ctx.ReplyText(err.Error())
					return
				}
				if service.Options.OnEnable != nil {
					service.Options.OnEnable(ctx)
				} else {
					ctx.ReplyText("启用成功")
				}
			case "禁用":
				if err := service.Disable(grp); err != nil {
					ctx.ReplyText(err.Error())
					return
				}
				if service.Options.OnDisable != nil {
					service.Options.OnDisable(ctx)
				} else {
					ctx.ReplyText("禁用成功")
				}
			}
		})

		// todo 启用、禁用全部插件在某个群或某个私聊
		robot.OnCommandGroup([]string{"启用全部", "禁用全部"}, robot.UserOrGroupAdmin).SetBlock(true).FirstPriority().Handle(func(ctx *robot.Ctx) {
		})

		// 启用、禁用某个插件在所有群和所有私聊
		robot.OnCommandGroup([]string{"全局启用", "全局禁用"}, robot.UserOrGroupAdmin).SetBlock(true).FirstPriority().Handle(func(ctx *robot.Ctx) {
			args := ctx.State["args"].(string)
			if args == "" {
				return
			}
			service, ok := managers.Lookup(args)
			if !ok {
				ctx.ReplyTextAndAt("没有找到对应插件服务")
				return
			}
			switch ctx.State["command"].(string) {
			case "全局启用":
				if service.Enable("all") != nil {
					ctx.ReplyText("全局启用失败")
					return
				}
				if service.Options.OnEnable != nil {
					service.Options.OnEnable(ctx)
				} else {
					ctx.ReplyText("全局启用成功")
				}
			case "全局禁用":
				if service.Disable("all") != nil {
					ctx.ReplyText("全局禁用失败")
					return
				}
				if service.Options.OnDisable != nil {
					service.Options.OnDisable(ctx)
				} else {
					ctx.ReplyText("全局禁用成功")
				}
			}
		})

		// 开启、关闭某个插件的全局模式
		robot.OnCommand("关闭全局模式", robot.UserOrGroupAdmin).SetBlock(true).FirstPriority().Handle(func(ctx *robot.Ctx) {
			args := ctx.State["args"].(string)
			if args == "" {
				return
			}
			service, ok := managers.Lookup(args)
			if !ok {
				ctx.ReplyTextAndAt("没有找到对应插件服务")
				return
			}
			if err := service.CloseGlobalMode(); err != nil {
				ctx.ReplyText(err.Error())
				return
			} else {
				ctx.ReplyText("关闭全局模式成功")
			}
		})
	})
}
