package manager

import (
	"github.com/imroc/req/v3"
	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/robot"
)

func registerCommand() {
	engine := control.Register("command", &control.Options[*robot.Ctx]{
		HideMenu: true,
	})

	// @机器人的命令
	engine.OnMessage(robot.OnlyAtMe).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		ctx.ReplyTextAndAt("您可以发送menu | 菜单解锁更多功能😎")
	})

	// 菜单输出
	engine.OnFullMatchGroup([]string{"menu", "菜单"}).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		options := control.GetOptionsMenu(ctx.Event.FromUniqueID)
		if options == nil || len(options.Menus) == 0 {
			ctx.ReplyTextAndAt("当前没有注册任何插件")
			return
		}

		// 🔔实现方案一：直接输出菜单
		//menus := "当前支持的功能有: \n"
		//for i := range options.Menus {
		//	menu := ""
		//	menu += "服务名: %s\n"
		//	menu += "别称: %s\n"
		//	menu += "默认开启状态: %v\n"
		//	menu += "当前开启状态: %v\n"
		//	menu += "插件描述: %s\n\n"
		//	menus += fmt.Sprintf(menu, options.Menus[i].Name, options.Menus[i].Alias, options.Menus[i].DefStatus, options.Menus[i].CurStatus, options.Menus[i].Describe)
		//}
		//ctx.ReplyTextAndAt(menus)

		// 🔔实现方案二：调用接口输出菜单（仅限作者个人使用，其他开发者请使用方案一或者自行修改）
		if err := req.C().Post("https://bot.yqqy.top/api/menu").SetBodyJsonMarshal(options).Do().Error(); err != nil {
			ctx.ReplyTextAndAt("菜单获取失败，请联系管理员")
			return
		}
		ctx.ReplyShareLink(robot.BotConfig.BotNickname, "机器人当前所有的指令都在这里哦！", "https://imgbed.link/file/10160", "https://bot.yqqy.top/menu?wxId="+ctx.Event.FromUniqueID)
	})
}
