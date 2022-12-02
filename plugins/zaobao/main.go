package zaobao

import (
	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/robot"
)

func init() {
	engine := control.Register("zaobao", &control.Options[*robot.Ctx]{
		Alias: "每日早报",
		Help:  "输入 {每日早报|早报} => 获取每天60s读懂世界",
	})

	engine.OnFullMatchGroup([]string{"早报", "每日早报"}).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		ctx.ReplyImage("https://api.qqsuu.cn/api/dm-60s?type=image")
	})
}
