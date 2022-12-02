package control

import (
	"sync"

	"github.com/yqchilde/wxbot/engine/robot"
)

var (
	// managers 每个插件对应的管理
	managers = NewManager[*robot.Ctx]("data/manager/plugins.db")
)

func newControl(service string, o *Options[*robot.Ctx]) robot.Rule {
	c := managers.NewControl(service, o)
	return func(ctx *robot.Ctx) bool {
		ctx.State["manager"] = c
		return c.Handler(ctx.Event.Message.FromGroup, ctx.Event.Message.FromWxid)
	}
}

func init() {
	once := sync.Once{}
	once.Do(func() {
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
			grp := ctx.Event.Message.FromGroup
			if grp == "" {
				// 个人用户
				grp = ctx.Event.Message.FromWxid
			}
			switch ctx.State["command"].(string) {
			case "启用":
				if service.Enable(grp) != nil {
					ctx.ReplyText("启用失败")
					return
				}
				if service.Options.OnEnable != nil {
					service.Options.OnEnable(ctx)
				} else {
					ctx.ReplyText("启用成功")
				}
			case "禁用":
				if service.Disable(grp) != nil {
					ctx.ReplyText("禁用失败")
					return
				}
				if service.Options.OnDisable != nil {
					service.Options.OnDisable(ctx)
				} else {
					ctx.ReplyText("禁用成功")
				}
			}
		})
	})
}
