package control

import (
	"fmt"
	"strings"
	"sync"

	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/robot"
)

var (
	once     = sync.Once{}
	managers = NewManager("data/manager/manager.db")
)

func newControl(service string, o *Options) robot.Rule {
	c := managers.NewControl(service, o)
	return func(ctx *robot.Ctx) bool {
		ctx.State["manager"] = c
		return c.Handler(ctx.Event.FromGroup, ctx.Event.FromWxId)
	}
}

func init() {
	robot.RegisterApi(&controlApi{})
	once.Do(func() {
		// 记录聊天文本消息
		robot.OnMessage().SetBlock(false).Handle(func(ctx *robot.Ctx) {
			if !ctx.IsText() {
				return
			}
			var msgType string
			if ctx.IsEventGroupChat() {
				msgType = "group"
			} else if ctx.IsEventPrivateChat() {
				msgType = "private"
			} else {
				return
			}

			if err := managers.D.Table("__message").Create(&MessageRecord{
				&robot.MessageRecord{
					Type:       msgType,
					FromWxId:   ctx.Event.FromUniqueID,
					FromNick:   ctx.Event.FromUniqueName,
					SenderWxId: ctx.Event.FromWxId,
					SenderNick: ctx.Event.FromName,
					Content:    ctx.MessageString(),
				},
			}).Error; err != nil {
				log.Errorf("记录消息失败: %v", err)
			}
		})

		// 响应或沉默某个群或某个私聊，沉默后在该群或私聊中的消息不会被机器人响应
		robot.OnCommandGroup([]string{"响应", "沉默"}, robot.UserOrGroupAdmin).SetBlock(true).FirstPriority().Handle(func(ctx *robot.Ctx) {
			args := ctx.State["args"].(string)
			if args == "" {
				args = ctx.Event.FromUniqueID
			}
			switch ctx.State["command"].(string) {
			case "响应":
				if err := managers.Response(args); err != nil {
					ctx.ReplyTextAndAt("ERROR: " + err.Error())
				} else {
					ctx.ReplyTextAndAt(fmt.Sprintf("开始响应[%s]消息啦~", args))
				}
			case "沉默":
				if err := managers.Silence(args); err != nil {
					ctx.ReplyTextAndAt("ERROR: " + err.Error())
				} else {
					ctx.ReplyTextAndAt(fmt.Sprintf("已经沉默[%s]消息啦!", args))
				}
			}
		})

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

		// 在某个群ban、unban某个用户
		robot.OnCommandGroup([]string{"ban", "unban"}, robot.UserOrGroupAdmin).SetBlock(true).FirstPriority().Handle(func(ctx *robot.Ctx) {
			args := strings.Split(ctx.State["args"].(string), " ")
			if len(args) == 0 {
				return
			}

			serv, wxId, grp := args[0], "", ctx.Event.FromUniqueID
			if ctx.IsReference() {
				wxId = ctx.Event.ReferenceMessage.ChatUser
			} else {
				wxId = args[1]
			}

			service, ok := managers.Lookup(serv)
			if !ok {
				ctx.ReplyTextAndAt("没有找到对应插件服务")
				return
			}
			switch ctx.State["command"].(string) {
			case "ban":
				if err := service.Ban(wxId, grp); err != nil {
					ctx.ReplyText(err.Error())
					return
				}
				ctx.ReplyText("ban成功")
			case "unban":
				if err := service.UnBan(wxId, grp); err != nil {
					ctx.ReplyText(err.Error())
					return
				}
				ctx.ReplyText("unban成功")
			}
		})
	})
}
