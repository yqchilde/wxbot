package invite

import (
	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/robot"
)

func init() {
	engine := control.Register("invite", &control.Options{
		Alias: "邀请好友",
	})

	// 邀请好友进群
	engine.OnFullMatch("测试发邀请链接").SetBlock(true).Handle(func(ctx *robot.Ctx) {
		// typ:1-直接拉，2-发送邀请链接
		ctx.InviteIntoGroup("39171925457@chatroom", "wxid", 2)
	})
}
