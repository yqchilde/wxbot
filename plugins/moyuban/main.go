package moyuban

import (
	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/robot"
)

func init() {
	engine := control.Register("moyu", &control.Options[*robot.Ctx]{
		Alias: "摸鱼日历",
		Help:  "输入 {摸鱼日历|摸鱼} => 获取摸鱼办日历",
	})
	engine.OnFullMatchGroup([]string{"摸鱼日历", "摸鱼"}).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		ctx.ReplyImage("https://api.vvhan.com/api/moyu")
	})
}
