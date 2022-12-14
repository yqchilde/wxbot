package control

import (
	"sync"

	"github.com/imroc/req/v3"

	"github.com/yqchilde/wxbot/engine/robot"
)

var (
	once = sync.Once{}
	// managers 每个插件对应的管理
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
		robot.On(robot.OnlyAtMe).SetBlock(true).Handle(func(ctx *robot.Ctx) {
			ctx.ReplyTextAndAt("您可以发送menu | 菜单解锁更多功能😎")
		})

		robot.OnFullMatchGroup([]string{"menu", "菜单"}).SetBlock(true).Handle(func(ctx *robot.Ctx) {
			services := managers.LookupAll()
			data := make(map[string]interface{})
			data["wxId"] = ctx.Event.FromUniqueID
			data["menus"] = make([]map[string]interface{}, 0, len(services))
			for _, s := range services {
				if !s.Options.ShowMenu {
					continue
				}
				data["menus"] = append(data["menus"].([]map[string]interface{}), map[string]interface{}{
					"name":      s.Service,
					"alias":     s.Options.Alias,
					"priority":  s.Options.priority,
					"describe":  s.Options.Help,
					"defStatus": !s.Options.DisableOnDefault,
					"curStatus": s.IsEnabledIn(ctx.Event.FromUniqueID),
				})
			}
			if err := req.C().Post("https://bot.yqqy.top/api/menu").SetBodyJsonMarshal(data).Do().Error(); err != nil {
				ctx.ReplyTextAndAt("菜单获取失败，请联系管理员")
				return
			}
			ctx.ReplyShareLink(robot.BotConfig.BotNickname, "机器人当前所有的指令都在这里哦！", "https://imgbed.link/file/10160", "https://bot.yqqy.top/menu?wxId="+ctx.Event.FromUniqueID)
		})

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
